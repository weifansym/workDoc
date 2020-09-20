# Best Practices for ES6 Promises
ES6 promises are great! They are integral constructs for asynchronous programming in JavaScript, ultimately replacing the old callback-based pattern that was 
most infamously known for bringing about deeply nested code ("callback hell").

Unfortunately, promises are not exactly the easiest concept to grasp. In this article, I will discuss the best practices I have learned over the years that 
helped me make the most out of asynchronous JavaScript.

Handle promise rejections
Nothing is more frustrating that an unhandled promise rejection. This occurs when a promise throws an error but no Promise#catch handler exists to gracefully 
handle it.

When debugging a heavily concurrent application, the offending promise is incredibly difficult to find due to the cryptic (and rather intimidating) error message
that follows. However, once it is found and deemed reproducible, the state of the application is often just as difficult to determine due to all the concurrency 
in the application itself. Overall, it is not a fun experience.

The solution, then, is simple: always attach a Promise#catch handler for promises that may reject, no matter how unlikely.

Besides, in future versions of Node.js, unhandled promise rejections will crash the Node process. There is no better time than now to make graceful error handling
a habit.

## Keep it "linear"
> [Please don't nest promises](https://dev.to/somedood/please-don-t-nest-promises-3o1o)

In a recent article, I explained why it is important to avoid nesting promises. In short, nested promises stray back into the territory of "callback hell". 
The goal of promises is to provide idiomatic standardized semantics for asynchronous programming. By nesting promises, we are vaguely returning to the verbose 
and rather cumbersome error-first callbacks popularized by Node.js APIs.

To keep asynchronous activity "linear", we can make use of either [asynchronous functions](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function) or properly chained promises.

```import { promises as fs } from 'fs';

// Nested Promises
fs.readFile('file.txt')
  .then(text1 => fs.readFile(text1)
    .then(text2 => fs.readFile(text2)
      .then(console.log)));

// Linear Chain of Promises
const readOptions = { encoding: 'utf8' };
const readNextFile = fname => fs.readFile(fname, readOptions);
fs.readFile('file.txt', readOptions)
  .then(readNextFile)
  .then(readNextFile)
  .then(console.log);

// Asynchronous Functions
async function readChainOfFiles() {
  const file1 = await readNextFile('file.txt');
  const file2 = await readNextFile(file1);
  console.log(file2);
}
```
## util.promisify is your best friend
As we transition from error-first callbacks to ES6 promises, we tend to develop the habit of "promisifying" everything.

For most cases, wrapping old callback-based APIs with the Promise constructor will suffice. A typical example is "promisifying" globalThis.setTimeout 
as a sleep function.
```
const sleep = ms => new Promise(
  resolve => setTimeout(resolve, ms)
);
await sleep(1000);
```
However, other external libraries may not necessarily "play nice" with promises out of the box. Certain unforeseen side effects—such as memory leaks—may 
occur if we are not careful. In Node.js environments, the util.promisify utility function exists to tackle this issue.

As its name suggests, util.promisify corrects and simplifies the wrapping of callback-based APIs. It assumes that the given function accepts an error-first 
callback as its final argument, as most Node.js APIs do. If there exists special implementation details1, library authors can also provide a ["custom promisifier"](https://nodejs.org/api/util.html#util_custom_promisified_functions).
```
import { promisify } from 'util';
const sleep = promisify(setTimeout);
await sleep(1000);
```
## Avoid the sequential trap
> [Avoiding the Sequential Trap](https://dev.to/somedood/javascript-concurrency-avoiding-the-sequential-trap-7f0)

In the previous article in this series, I extensively discussed the power of scheduling multiple independent promises. Promise chains can only get us so far when 
it comes to efficiency due to its sequential nature. Therefore, the key to minimizing a program's "idle time" is concurrency.
```
import { promisify } from 'util';
const sleep = promisify(setTimeout);

// Sequential Code (~3.0s)
sleep(1000)
  .then(() => sleep(1000));
  .then(() => sleep(1000));

// Concurrent Code (~1.0s)
Promise.all([ sleep(1000), sleep(1000), sleep(1000) ]);
```
Beware: promises can also block the event loop
Perhaps the most popular misconception about promises is the belief that promises allow the execution of "multi-threaded" JavaScript. Although the event loop 
gives the illusion of "parallelism", it is only that: an illusion. Under the hood, JavaScript is still single-threaded.

The event loop only enables the runtime to concurrently schedule, orchestrate, and handle events throughout the program. Loosely speaking, these "events" indeed 
occur in parallel, but they are still handled sequentially when the time comes.

In the following example, the promise does not spawn a new thread with the given executor function. In fact, the executor function is always executed immediately 
upon the construction of the promise, thus blocking the event loop. Once the executor function returns, top-level execution resumes. Consumption of the resolved 
value (through the Promise#then handler) is deferred until the current call stack finishes executing the remaining top-level code.2
```
console.log('Before the Executor');

// Blocking the event loop...
const p1 = new Promise(resolve => {
  // Very expensive CPU operation here...
  for (let i = 0; i < 1e9; ++i)
    continue;
  console.log('During the Executor');
  resolve('Resolved');
});

console.log('After the Executor');
p1.then(console.log);
console.log('End of Top-level Code');

// Result:
// 'Before the Executor'
// 'During the Executor'
// 'After the Executor'
// 'End of Top-level Code'
// 'Resolved'
```
Since promises do not automatically spawn new threads, CPU-intensive work in subsequent Promise#then handlers also blocks the event loop.
```
Promise.resolve()
//.then(...)
//.then(...)
  .then(() => {
    for (let i = 0; i < 1e9; ++i)
      continue;
  });
```
Take memory usage into consideration
Due to some unfortunately necessary [heap allocations](https://www.youtube.com/watch?v=wJ1L2nSIV1s), promises tend to exhibit relatively hefty memory footprints and computational costs.

In addition to storing information about the Promise instance itself (such as its properties and methods), the JavaScript runtime also dynamically allocates 
more memory to keep track of the asynchronous activity associated with each promise.

Furthermore, given the Promise API's extensive use of closures and callback functions (both of which require heap allocations of their own), a single promise 
surprisingly entails a considerable amount of memory. An array of promises can prove to be quite consequential in hot code paths.

As a general rule of thumb, each new instance of a Promise requires its own hefty heap allocation for storing properties, methods, closures, and asynchronous 
state. The less promises we use, the better off we'll be in the long run.

## Synchronously settled promises are redundant and unnecessary
As discussed earlier, promises do not magically spawn new threads. Therefore, a completely synchronous executor function (for the Promise constructor) only has 
the effect of introducing an unnecessary layer of indirection.3
```
const promise1 = new Promise(resolve => {
  // Do some synchronous stuff here...
  resolve('Presto');
});
```
Similarly, attaching Promise#then handlers to synchronously resolved promises only has the effect of slightly deferring the execution of code.4 For this use case,
one would be better off using global.setImmediate instead.
```
promise1.then(name => {
  // This handler has been deferred. If this
  // is intentional, one would be better off
  // using `setImmediate`.
});
```
Case in point, if the executor function contains no asynchronous I/O operations, it only serves as an unnecessary layer of indirection that bears the 
aforementioned memory and computational overhead.

For this reason, I personally discourage myself from using Promise.resolve and Promise.reject in my projects. The main purpose of these static methods is to 
optimally wrap a value in a promise. Given that the resulting promise is immediately settled, one can argue that there is no need for a promise in the first 
place (unless for the sake of API compatibility).
```
// Chain of Immediately Settled Promises
const resolveSync = Promise.resolve.bind(Promise);
Promise.resolve('Presto')
  .then(resolveSync)  // Each invocation of `resolveSync` (which is an alias
  .then(resolveSync)  // for `Promise.resolve`) constructs a new promise
  .then(resolveSync); // in addition to that returned by `Promise#then`.
```
### Long promise chains should raise some eyebrows
There are times when multiple asynchronous operations need to be executed in series. In such cases, promise chains are the ideal abstraction for the job.

However, it must be noted that since the Promise API is meant to be chainable, each invocation of Promise#then constructs and returns a whole new Promise 
instance (with some of the previous state carried over). Considering the additional promises constructed by intermediate handlers, long chains have the potential 
to take a significant toll on both memory and CPU usage.
```
const p1 = Promise.resolve('Presto');
const p2 = p1.then(x => x);

// The two `Promise` instances are different.
p1 === p2; // false
```
Whenever possible, promise chains must be kept short. An effective strategy to enforce this rule is to disallow fully synchronous Promise#then handlers except 
for the final handler in the chain.

In other words, all intermediate handlers must strictly be asynchronous—that is to say, they return promises. Only the final handler reserves the right to run 
fully synchronous code.
```
import { promises as fs } from 'fs';

// This is **not** an optimal chain of promises
// based on the criteria above.
const readOptions = { encoding: 'utf8' };
fs.readFile('file.txt', readOptions)
  .then(text => {
    // Intermediate handlers must return promises.
    const filename = `${text}.docx`;
    return fs.readFile(filename, readOptions);
  })
  .then(contents => {
    // This handler is fully synchronous. It does not
    // schedule any asynchronous operations. It simply
    // processes the result of the preceding promise
    // only to be wrapped (as a new promise) and later
    // unwrapped (by the succeeding handler).
    const parsedInteger = parseInt(contents);
    return parsedInteger;
  })
  .then(parsed => {
    // Do some synchronous tasks with the parsed contents...
  });
```
As demonstrated by the example above, fully synchronous intermediate handlers bring about the redundant wrapping and unwrapping of promises. This is why it is 
important to enforce an optimal chaining strategy. To eliminate redundancy, we can simply integrate the work of the offending intermediate handler into the 
succeeding handler.
```
import { promises as fs } from 'fs';

const readOptions = { encoding: 'utf8' };
fs.readFile('file.txt', readOptions)
  .then(text => {
    // Intermediate handlers must return promises.
    const filename = `${text}.docx`;
    return fs.readFile(filename, readOptions);
  })
  .then(contents => {
    // This no longer requires the intermediate handler.
    const parsed = parseInt(contents);
    // Do some synchronous tasks with the parsed contents...
  });
```
Keep it simple!
If you don't need them, don't use them. It's as simple as that. If it's possible to implement an abstraction without promises, then we should always prefer 
that route.

Promises are not "free". They do not facilitate "parallelism" in JavaScript by themselves. They are simply a standardized abstraction for scheduling and handling 
asynchronous operations. If the code we write isn't inherently asynchronous, then there is no need for promises.

Unfortunately, more often than not, we do need promises for powerful applications. This is why we have to be cognizant of all the best practices, trade-offs, 
pitfalls, and misconceptions. At this point, it is only a matter of minimizing usage—not because promises are "evil", but because they are so easy to misuse.

But this is not where the story ends. In the [next part](https://dev.to/somedood/best-practices-for-es2017-asynchronous-functions-async-await-39ji) of this series, I will extend the discussion of best practices to [ES2017 asynchronous functions](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Asynchronous/Async_await) 
(async/await).

1. This may include specific argument formats, initialization operations, clean-up operations, and so on and so forth.
2. In essence, this is what it means to schedule a "microtask" in the "microtask queue". Once the current top-level code finishes executing, 
the "microtask queue" waits for all scheduled promises to be settled. Over time, for each resolved promise, the "microtask queue" invokes the respective 
Promise#then handler with the resolved value (as stored by the resolve callback). 
3. With the added overhead of a single promise. 
4. With the added overhead of constructing a new promise for each chained handler. 

转自：https://dev.to/somedood/best-practices-for-es6-promises-36da







