## Learning the Node.js Runtime
[原地址](https://jscomplete.com/learn/node-beyond-basics)
I believe the majority of developers learn Node the wrong way. Most tutorials, books, and courses about Node focus on the Node
ecosystem – not the Node runtime itself. They focus on teaching what can be done with all the packages available for you when 
you work with Node, like Express and Socket.IO, rather than teaching the capabilities of the Node runtime itself.

There are good reasons for this. Node is raw and flexible. It doesn’t provide complete solutions, but rather provides a rich 
runtime that enables you to implement solutions of your own. Libraries like Express.js and Socket.IO are more of complete 
solutions, so it makes more sense to teach those libraries, so you can enable learners to use these complete solutions.

The conventional wisdom seems to be that only those whose job is to write libraries like Express.js and Socket.IO need to 
understand everything about the Node runtime, but I think this is wrong. A solid understanding of the Node runtime itself is 
the best thing you can do before using those complete solutions. You should at least have the knowledge and confidence to
judge a package by its code so you can make an educated decision about using it.
### The Node Knowledge Challenge

Let me give you a taste of the kind of questions you will be able to answer after reading this book. Look at this as your Node
knowledge challenge. If you can answer most of these questions, this book is probably not for you.

* What is the relationship between Node and V8? Can Node work without V8?

* come when you declare a global variable in any Node file it’s not really global to all modules?

* When exporting the API of a Node module, why can we sometimes use exports and other times we have to use module.exports?

* What is the Call Stack? Is it part of V8?

* What is the Event Loop? Is it part of V8?

* What is the difference between setImmediate and process.nextTick?

* What are the major differences between spawn, exec, and fork?

* How does the cluster module work? How is it different than using a load balancer?

* What will Node do when both the call stack and the event loop queue are empty?

* What are V8 object and function templates?

* What is libuv and how does Node use it?

* How can we do one final operation before a Node process exits? Can that operation be done asynchronously?

* Besides V8 and libuv, what other external dependencies does Node have?

* What’s the problem with the process uncaughtException event? How is it different than the exit event?

* What are the 5 major steps that the require function does?

* How can you check for the existence of a local module?

* What are circular modular dependencies in Node and how can they be avoided?

* What are the 3 file extensions that will be automatically tried by the require function?

* When creating an http server and writing a response for a request, why is the end() function required?

* When is it ok to use the file system *Sync methods?

* How can you print only one level of a deeply nested object?

* How come top-level variables are not global?

* The objects exports, require, and module are all globally available in every module but they are different in every module.
How?

* If you execute a JavaScript file that has the single line: console.log(arguments); with Node, what exactly will Node print?

* How can a module be both "requirable" by other modules and executable directly using the node command?

* What’s an example of a built-in stream in Node that is both readable and writable?

* What happens when the line cluster.fork() gets executed in a Node script?

* What’s the difference between using event emitters and using simple callback functions to allow for asynchronous handling of code?

* What’s the difference between the Paused and the Flowing modes of readable streams?

* How can you read data from a connected socket?

* The require function always caches the module it requires. What can you do if you need to execute the code in a required 
module many times?

* When working with streams, when do you use the pipe function and when do you use events? Can those two methods be combined?

### Fundamentals
Okay, I would categorize some of the questions above as fundamentals. Let me start by answering these:

