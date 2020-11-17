## Investigating Multithreaded Implementation of rstm2d

Here are some benchmarks for my quick and dirty multithreaded implementation of rsmt2d, which can be found on the 'bench' branch of [my fork of rsmt2d](https://github.com/evan-forbes/rsmt2d/tree/bench)

![Performance](performance.png) 
Overall, the performance is increased roughly 4 fold on an 8 core cpu. 
![Overhead](overhead.png)
The overhead was determined by comparing the current implementation with the multithreaded version limited to a single thread. While the overhead is within the margin of error time wise, I expect other measurements of overhead to be well above the current implementation due to the extra allocations required.