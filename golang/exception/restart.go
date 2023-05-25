package main


func main() {
	ctx, cancel := context.WithCancel(context.Background())

	_ = pool.GetPoolInstance().Submit(func() {
		RunCpuPprof(ctx)
	})

	_ = pool.GetPoolInstance().Submit(func() {
		RewriteStderrFile()
	})

	initRuntime(ctx)
	defer runtimeClose(ctx)

	initServer(ctx, cancel)
}

// --- for performance tracking --- //
func RunCpuPprof(ctx context.Context) {
	_ = os.Mkdir("pprof", os.ModePerm)
	for {
		func() {
			fileName := fmt.Sprintf("./pprof/%s.pprof", time.Now().Format("2006-01-02 15:04:05"))
			cpuFile, err := os.Create(fileName)
			if err != nil {
				xlog.Errorf("[RunCpuPprof] err: %s", err.Error())
			}
			defer cpuFile.Close()
			_ = pprof.StartCPUProfile(cpuFile)
			defer pprof.StopCPUProfile()
			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Minute):
			}
		}()
	}
}

const stdErrFile = "/tmp/go-app1-stderr.log"

var stdErrFileHandler *os.File

func RewriteStderrFile() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return nil
}

// --- /for performance tracking --- //

// --- for runtime --- //
func initRuntime(ctx context.Context) {
	initConfig()
	initLog()
	initRepository(ctx)
}
func initConfig() {
	xlog.Info("[initConfig] start.")
	config.Init()
	xlog.Info("[initConfig] finished.")
}
func initLog() {
	xlog.Infof("[initLog] start.")
	logconf := xlog.NewProductionConfig()
	logconf.Level.SetLevel(zapcore.Level(config.Conf.Log.Level))
	logconf.Name = config.Conf.App.Name
	logconf.Writer = os.Stdout
	logconf.Encoding = encoder.LIVE
	l, err := xlog.New(logconf)
	if err != nil {
		xlog.Errorf("[initLog] err: %s", err.Error())
		panic(err)
	}
	xlog.UpdateGlobalLog(l)
	xlog.Infof("[initLog] finished.")
}
func initRepository(ctx context.Context) {
	xlog.Infof("[initRepository] start.")
	repository.Init(
		repository.WithMySQL(ctx),
		repository.WithRedis(ctx),
		repository.WithMongo(ctx),
	)
	xlog.Infof("[initRepository] finished.")
}
func runtimeClose(ctx context.Context) {
	repository.Close(ctx)
}

// --- /for runtime --- //

// --- servers --- //
func initServer(ctx context.Context, cancel context.CancelFunc) {
	errCh := make(chan error)

	// create http server
	httpCloseCh := make(chan struct{})
	StartHttpServer(ctx, errCh, httpCloseCh)
	// create rpc server

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		cancel()
		xlog.Info("[initServer] Start graceful shutdown")
	case err := <-errCh:
		cancel()
		xlog.Errorf("[initServer] server err: %s", err.Error())
	}

	// wait for http server
	<-httpCloseCh
	xlog.Infof("[initServer] http server %s exit!", config.Conf.App.Name)
}
func initNewStructure(engine *gin.Engine) {
	// 初始化 access logger
	var log xlog2.Logger
	if env.Active().IsLocal() {
		log = xlog2.NewJSONLogger(xlog2.WithFileRotationP(config.Conf.Log.FilePath))
	} else {
		log = xlog2.NewJSONLogger()
	}
	router.NewHTTPServer(engine, log)
}
func StartHttpServer(ctx context.Context, errChan chan error, httpCloseCh chan struct{}) {
	engine := gin.Default()
	initNewStructure(engine)

	// init server
	srv := &http.Server{
		Addr:    config.Conf.HttpServer.Addr,
		Handler: engine,
	}

	// run http server
	_ = pool.GetPoolInstance().Submit(func() {
		xlog.Info(fmt.Sprintf("%s server is starting on %s", config.Conf.App.Name, config.Conf.HttpServer.Addr))
		errChan <- srv.ListenAndServe()
	})

	// run rpc server
	_ = pool.GetPoolInstance().Submit(func() {
		// keep listening targeted port...
		listener, err := net.Listen("tcp", config.Conf.RPCServer.Addr)
		if err != nil {
			log.Fatal("[rpc] ListenTCP error:", err)
		}
		defer listener.Close()

		for {
			// got connection
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("[rpc] Accept error:", err)
			}

			// create serve connection
			rpc.ServeConn(conn)
		}
	})

	// watch the ctx exit
	_ = pool.GetPoolInstance().Submit(func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			xlog.Info(fmt.Sprintf("httpServer shutdown:%v", err))
		}
		httpCloseCh <- struct{}{}
	})
}

// --- /servers --- //
