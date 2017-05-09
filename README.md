### Tuenti Challenge 7. Problem 11

This repo solves the 11th problem of the [Tuenti Challenge](https://contest.tuenti.net/Challenges?id=11) from 2017.

#### How to
To build the project just do `go build *.go`. To use the program pipe the input to main, the solution is outputed to the stdin:  `less input | ./main > output`.

#### The problem

The full description is in the [link](https://contest.tuenti.net/Challenges?id=11), it is a bit too long but it can be synthesized as follows. There is a directed graph we want to traverse, for the problem we need to find the time required to visit each node. Contrary to a usual traversing problem (say Dijkstra) the edges don't take time to be travelled, but they cost some "energy" which is recharged at the nodes. This energy comes in several colours and it takes some time to recharge.

#### The solution

So, the first thing was to model the colors. There is a big hint in the statement, which suggests using binary code:

   > For example, if you are Red and you absorb Blue energy you would become Purple. Note that if you are Purple and you absorb Blue again, nothing would happen
  
The way to proceed is as follows, each primary color will be a power of two. If red was the first primary color then red is represented by 1, blue is the second so 2... In general it will be `1 << n`. Secondary colors are obtained by binary `or`. For example: 

    purple := blue | red  = 1 | 2 = 3

Once we have the colours, we have to plan how we travel the graph, Dijkstra is a clear candidate, at each step we consider all the options and we only plan those actions that could be beneficial, we plug them into a priority queue and pull the next actions from the queue.

There are a couple of bits missing to be able to apply Dijkstra. First is that it normally only refers to connecting two points, but that is not a problem we just keep going and then return the distance to all the points. Second and most important is that normally the only thing that matters is the distance, but in this problem we are also worry about the energy. For example, we might arrive to some node in 6 time units with red energy and in 9 time units with purple energy. Normally we would discard the 9 but in this case it might be beneficial to reach further nodes.
 
 Instead of saving the minimum distance to reach each vertex, for each color we store the minimum distance that we needed to arrive to that vertex with that color. So in the previous case we would have `{1: 6, 3: 9}` because red had code 1 and purple code 3. If later on we arrive to the same node with color blue and 10 units we would discard the action because blue is contained in purple and we already got it in a faster way.
 
 
 Apart from that, there are some minor optimizations like precomputing all possible charges so that we don't do two charges in a row and similar things.
 
 About the language, in the previous challenges I chose node and I ended up having some problems of performance, so this time I gave it a try to golang. And I was surprised beyond my expectations, the language was nice to write and it didn't event sweat, less than a minute for the whole algorithm.
 
The only pity is that I assumed that all secondary colors where defined from primary ones, but that was not the case in the final examples. For example in the 370th `Vbdm` was composed of `Djnx`,  `Xbovbt`, `Uhfcps` and `Hexjegb` none of which was primary. This is already solved in the repo.

  
 