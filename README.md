## Investigating Multithreaded Implementation of rstm2d

Here are some benchmarks for my quick and dirty multithreaded implementation of rsmt2d, which can be found on the 'bench' branch of [my fork of rsmt2d](https://github.com/evan-forbes/rsmt2d/tree/bench)

![Performance](performance.png) 
The performance is increased roughly 4 fold on an 8 core cpu. 
![Overhead](overhead.png)
The overhead was determined by comparing the current implementation with the multithreaded version limited to a single thread. While the overhead is within the margin of error time wise, I expect other measurements of overhead to be well above the current implementation due to the extra allocations required.

![trace](multithread_8_trace.png)
The trace is... nasty. There seems to be quite of lot of room for improvement, but I'm not sure how much due to the [3rd quadrant](https://github.com/lazyledger/lazyledger-specs/blob/master/specs/figures/rs2d_quadrants.svg) needing to be computed after the 2nd.

spreadsheet [here](https://docs.google.com/spreadsheets/d/1oLfHhEMRSsz99A26wBLddiLgaZiJN0P9c4y1hoHj2IE/edit?usp=sharing)