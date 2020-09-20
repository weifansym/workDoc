## Introduction
Asynchronous functions and callback functions are at the heart of JavaScript's "single-threaded" concurrency model. When we talk about asynchronous operations 
in JavaScript, we often hear about the ingenious engineering behind the humble but legendary [event loop](https://www.youtube.com/watch?v=8aGhZQkoFbQ). Such discussions will immediately be followed by the 
obligatory statement of: "Never block the event loop."

Indeed, it is a "mortal sin" to block the event loop. For that matter, the event loop (of any program) is quite like the human heartbeat. If the heart continues 
to beat at a steady pace, the program runs smoothly. However, if certain blockages disturb the natural rhythm, then everything starts to break down.


## Scope and Limitations
In this series of articles, we will explore the various ways to optimize the execution of asynchronous operations, but not the operations themselves. 
This distinction must be made because optimizing the "operations themselves" implies the discussion of implementation-specific details and logic, which are 
beyond the scope of this article.

Instead, we will focus on the proper scheduling of such operations. As much as possible, the goal is to take advantage of concurrency whenever it is possible. 
The sequential execution of asynchronous operations is fine—or even necessary—in some cases, but to make the most out of asynchronous JavaScript, we must minimize 
the "idle" moments of a program.

## Idle Execution
A JavaScript program is considered to be "idle" when there is literally nothing blocking the event loop, yet the program continues to wait for pending asynchronous
operations. In other words, an "idle program" occurs when there is nothing left to do but wait. Let us consider the following example:

> DISCLAIMER: [Top-level](https://gist.github.com/Rich-Harris/0b6f317657f5167663b493c722647221) await will be used throughout this article for the sake of brevity and demonstration purposes. However, 
as [Rich Harris](https://github.com/Rich-Harris) (creator of [Svelte](https://github.com/sveltejs/svelte) and [Rollup](https://github.com/rollup/rollup)) argued, this would be a "footgun" to the reliable execution of top-level JavaScript in production environments. 
As such, it is not recommended to haphazardly use top-level await such as in the examples to follow.

```
// Assuming that this network request takes one second to respond...
await fetch('https://example.com');

// Anything after this point is code that cannot be
// executed until the network request resolves.
console.log('This will run one second later.'):
```
The issue with waiting for asynchronous code to finish is the "idle time" during which many other asynchronous operations could have been scheduled.

Alternatively, numerous synchronous computations could have also been scheduled in the meantime (via [worker threads](https://nodejs.org/api/worker_threads.html#worker_threads_worker_threads) and [web workers](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers), for example) so that when 
the network request finally finishes, everything is ready, set, computed, and cached by then.

Of course, if the forthcoming computations depend on the result of the network request, then it is totally necessary to wait. In such situations where 
asynchronous operations are meant to be executed sequentially, an effort must still be made to cut down on the program's "idle time". To demonstrate this, 
let us consider an example with the file system involved:
```
import fetch from 'node-fetch';
import { promises as fs } from 'fs';
import { promisify } from 'util';

const sleep = promisify(setTimeout);

async function purelySequential() {
  // Let us assume that this file contains a single line
  // of text that happens to be some valid URL.
  const url = await fs.readFile('file.txt');
  const response = await fetch(url);

  // Execute some **unrelated** asynchronous
  // opeartion here...
  await sleep(2500);

  return result;
}
```
The function above reads from a file and then uses the retrieved text as the URL input for a network request. Once the request resolves, it executes another 
asynchronous operation that takes at least 2.5 seconds to finish.

If all goes well, the minimum total execution time of the function is 2.5 seconds. Anything less than that is impossible because of the sequential nature of 
the function. It must first wait for the file read to finish before initializing the network request. Since we must await the fetch request, the execution of 
the function pauses until the Promise settles. All of these asynchronous operations must resolve before we can even schedule the unrelated asynchronous operation.

We can optimize this function by scheduling the latter operation while waiting for the file read and the network request to finish. However, it must be 
reiterated that this only works with the assumption that the latter operation does not depend on the output of the aforementioned asynchronous operations.
```
import fetch from 'node-fetch';
import { promises as fs } from 'fs';
import { promisify } from 'util';

const sleep = promisify(setTimeout);

async function optimizedVersion() {
  // Schedule the unrelated operation here. The removal of the
  // `await` keyword tells JavaScript that the rest of the code can
  // be executed without having to _wait_ for `operation` to resolve.
  const operation = sleep(2500);

  // Now that `operation` has been scheduled, we can
  // now initiate the file read and the network request.
  const url = await fs.readFile('file.txt');
  const result = await fetch(url);

  // Once the network request resolves, we can now wait for
  // the pending `operation` to resolve.
  await operation;

  return result;
}
```
Assuming that the file system and the network interactions are fast, the optimized function now has a maximum execution time of 2.5 seconds. This is good news! 
By cleverly scheduling asynchronous operations, we have optimized the code to run concurrently.

To truly drive this point home, the example below demonstrates the discussed pattern with the sleep utility function:
```
import { promisify } from 'util';
const sleep = promisify(setTimeout);

console.time('Sequential');
await sleep(1000);
await sleep(2000);
console.timeEnd('Sequential');

console.time('Optimized');
const operation = sleep(2000);
await sleep(1000);
await operation;
console.timeEnd('Optimized');

// Sequential: ~3.0 seconds ❌
// Optimized: ~2.0 seconds ✔
```
## Promise.all
In situations where multiple asynchronous operations are not required to run sequentially, we can make full use of JavaScript's concurrency model with Promise.all.
As a quick refresher, Promise.all accepts an array of promises and then returns a single promise that wraps the given array. Once all of the promises in the 
original array successfully resolve, Promise.all resolves with an array of the results.
```
const promises = [
  Promise.resolve(1),
  Promise.resolve(2),
  Promise.resolve(3),
];
const results = Promise.all(promises);

// [ 1, 2, 3 ]
console.log(await results);
```
Assuming that all promises are guaranteed to resolve, this presents us with the unique advantage of scheduling an array of concurrent promises. Let us consider 
the following example:
```
/**
 * This function runs three independent operations sequentially.
 * Even if each operation is independent from each other, it makes
 * the mistake of running one after the other as if they were
 * dependent. In this case, the "idle time" is unnecessary and
 * extremely wasteful.
 */
async function sequential() {
  await sleep(2000);
  await sleep(3000);
  await sleep(4000);
}

/**
 * This function runs all of the operations concurrently.
 * `Promise.all` automatically schedules all of the
 * promises in the given array. By the time they all
 * resolve, `Promise.all` can safely return the array
 * of resolved values (if applicable).
 */
async function concurrent() {
  await Promise.all([
    sleep(2000),
    sleep(3000),
    sleep(4000),
  ]);
}

// **TOTAL EXECUTION TIMES**
// Sequential: ~9.0 seconds ❌
// Concurrent: ~4.0 seconds ✔
```
## Promise.allSettled
However, there are times when we cannot assume the success of promises. More often than not, we have to handle errors. During those times, 
the new Promise.allSettled comes to the rescue.

As its name suggests, Promise.allSettled behaves in a similar manner to Promise.all. The main difference between the two is how they handle promise rejections. 
For Promise.all, if any of the promises in the input array fails, it will immediately terminate further execution and throw the rejected promise regardless 
of whether some promises were successful.
```
const results = Promise.all([
  Promise.resolve(1),
  Promise.reject(2),
  Promise.resolve(3),
]);

// 2
console.error(await results);
```
The issue with this approach is its "fail-fast" feature. What if we still want to retrieve the values of the resolved promises despite the error? That's exactly 
where Promise.allSettled shines. Instead of "failing fast", Promise.allSettled segregates the resolved promises from the rejected ones by marking them as 
either 'fulfilled' or 'rejected'. That way, we can work with the resolved values while still being able to handle any errors.
```
const results = Promise.allSettled([
  Promise.resolve(1),
  Promise.resolve(2),
  Promise.reject(3),
]);

// [
//   { status: 'fulfilled', value: 1 },
//   { status: 'fulfilled', value: 2 },
//   { status: 'rejected', reason: 3 },
// ]
console.log(await results);
```
## The Caveats of a Single-Threaded Language
Throughout the article, I have been very careful with the term "asynchronous operation". When first hearing about the wonders of ES6 promises, many JavaScript 
developers—myself included—have fallen victim to the misconception that JavaScript was suddenly "multi-threaded".

Promises enabled us to concurrently run multiple asynchronous operations, hence the illusion of "parallel execution". But alas, "free parallelism" could not be 
farther from the truth.

## I/O Operations
In JavaScript, it is important to differentiate input-output (I/O) operations from CPU-intensive tasks. An I/O operation—such as network and file system 
interactions—requires the program to wait until the data is ready to be consumed. However, this does not necessarily "block" the execution of the program. 
While waiting for an I/O operation to finish, the program can still execute other code. Optionally, the program can block itself and poll for the data.

For example, a program may ask the operating system to read a certain file. The operating system commands the hard drive to "spin some disks" and "flip some bits" 
until the file is completely read. Meanwhile, the program continues execution and calculates the digits of pi. Once the file is available, the program consumes 
the data.

With this example in mind, this is why I have also been careful with the word "scheduling". Asynchronous operations in JavaScript typically mean I/O operations 
and timeouts. When we fetch for a resource, we schedule a request and wait for the data to be available. Once the request is scheduled, we let the operating 
system "do its thing" so that other code in the program can execute for the meantime, hence Node.js' core tenet of "non-blocking I/O".

## CPU-Intensive Tasks
On the other hand, CPU-intensive tasks literally block the execution of a program due to expensive computations. This typically means lengthy search algorithms,
sort algorithms, regular expression evaluation, text parsing, compression, cryptography, and all sorts of math calculations.

In some cases, I/O operations can also block a program. However, that is usually a conscious design choice. Through the *-Sync functions, Node.js provides 
synchronous alternatives to certain I/O operations. Nonetheless, these synchronous activities are a necessary expense.

However, therein lies the issue: synchronicity is necessary. To work around this, the greatest minds in computer science introduced the notion of "multi-threaded
systems" in which code can run in parallel. By offloading computational work across multiple threads, computers became more efficient with CPU-intensive tasks.

Despite the potential of multi-threading, JavaScript was explicitly designed to be single-threaded simply because it was incredibly difficult to write "safe" 
and "correct" multi-threaded code. For the Web, this was a reasonable trade-off for the sake of security and reliability.

## Misconceptions with Promises
When ES6 promises came along, it was incredibly tempting to "promisify" everything. Promises gave the illusion that JavaScript was "multi-threaded" in some way.
A JavaScript runtime (such as Node.js and the browser) is indeed multi-threaded, but unfortunately, that does not mean JavaScript (the language) became anything
more than "single-threaded" per se.

In reality, promises still executed code synchronously, albeit at a later time. Contrary to intuition and idealisms, offloading a CPU-intensive task onto a
promise does not magically spawn a new thread. The purpose of a promise is not to facilitate parallelism, but to defer the execution of code until some data 
is resolved or rejected.

The key word here is "defer". By deferring execution, any computationally expensive task will still inevitably block the execution of a program—provided that 
the data is ready to be consumed by then.
```
// This promise will still block the event loop.
// It will **not** execute this in parallel.
new Promise(resolve => {
  calculateDigitsOfPi();
  mineForBitcoins();
  renderSomeGraphcs();
  doSomeMoreMath();
  readFileSync('file.txt');

  resolve();
});
```
## Promises and Worker Threads
As mentioned earlier, the main use case for promises is to defer the execution of code until the requested data is ready to be consumed. A promise implies the 
scheduling of an asynchronous I/O operation that will eventually resolve, but it does not mean parallelism for CPU-intensive tasks.

If parallelism for CPU-intensive tasks is absolutely necessary for an application, the best approach is to use [web workers](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers) in the browser. In Node.js,
[worker threads](https://nodejs.org/api/worker_threads.html#worker_threads_worker_threads) are the equivalent API.
If concurrency for multiple asynchronous I/O operations and timeouts is needed, promises and events are the best tools for the job.

When used incorrectly, a CPU-intensive task in a promise will block the event loop. Inversely, spreading multiple I/O operations across many background worker
threads is redundant and wasteful. By manually spawning a whole new thread just for an I/O operation, the thread is literally idle for most of its existence 
until the requested data arrives.

Delving into the more technical part of implementation details, a well-designed JavaScript runtime already handles and abstracts away [the multi-threaded aspect
of I/O operations](https://nodejs.org/en/docs/guides/dont-block-the-event-loop/). This is what makes the aforementioned misuse of worker threads "redundant".
Moreover, in Node.js, each background thread occupies a single slot in the worker pool. Given the fact that the number of threads in the [worker pool](https://nodejs.org/en/docs/guides/dont-block-the-event-loop/#don-t-block-the-worker-pool) is finite 
and limited, efficient thread management is critical to Node.js' ability to operate concurrently. Otherwise, redundantly spawning worker threads gravely 
mishandles the limited worker pool.

For this reason, an idle worker thread (due to pending I/O operations) is not only wasteful, but also unnecessary. One would be better off letting the JavaScript
runtime "do its thing" when handling I/O.

For this reason, an idle worker thread (due to pending I/O operations) is not only wasteful, but also unnecessary. One would be better off letting the JavaScript 
runtime "do its thing" when handling I/O.

## Conclusion
If there is one lesson to be learned from this article, it is the difference between I/O operations and CPU-intensive tasks. By understanding their use cases,
one can confidently identify the correct tools for maximizing JavaScript concurrency.

I/O operations inherently defer code until some data is ready. For this reason, we must make use of promises, events, and callback patterns to schedule requests.
With the proper orchestration of I/O operations, we can keep the event loop running while still being able to handle asynchronous code concurrently.

On the other hand, CPU-intensive tasks will inevitably block the execution of a program. Wisely offloading these synchronous operations to separate background 
worker threads is a surefire way of achieving parallelism. However, it is still of utmost importance to be cognisant of the overhead and the hidden costs of 
occupying a slot in the worker pool.

As a general rule of thumb, promises are for I/O operations while worker threads are for CPU-intensive tasks. Taking advantage of these core concepts helps us 
avoid the trap of sequential "blocking" code.

转自：https://dev.to/somedood/javascript-concurrency-avoiding-the-sequential-trap-7f0






