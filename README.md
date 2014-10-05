gochan-life
===========

Game of life implemented with goroutines and channels

## Why?

Because I wanted to see what that would look like.

## So, how deos this differ from every other implementation out there?

Channels drive all the interactions.  Each node in the game is a goroutine and
it is connected to its neighbors via channels.

The main engine function looks like this:

    func Run (state int,
              nw, n, ne, e, se, s, sw, w <-chan int,
              NW, N, NE, E, SE, S, SW, W chan<- int,
              Out chan<- Status) {

      // Start computation.
      for {
        // First, write our status
        NW <- state
        N  <- state
        NE <- state
        E  <- state
        SE <- state
        S  <- state
        SW <- state
        W  <- state

        // Count the state values coming in
        count := <-nw + <-n + <-ne + <-e + <-se + <-s + <-sw + <-w

        // Report to output
        Out <- Status { state > 0, count }

        // Determine our new state
        if state == 0 {
          if count == 3 {
            state = 1
          }
        } else {
          if count != 2 && count != 3 {
            state = 0
          }
        }
      }
    }


## So, how did it perform?

On my 1GB single core (i7) VM:

| Columns | Rows | Nodes | FPS | %RAM (via top)
|---------|------|-------|-----|-------------
| 78      | 20   | 1560  | 488 | 1.2%
| 160     | 20   | 3200  | 193 | 2.6%
| 160     | 50   | 8000  | 69  | 6.2%
| 160     | 100  | 16000 | 33  | 11.1%
| 160     | 140  | 22400 | 22  | ~24%
| 500     | 100  | 50000 | 6   | ~45%

Beyond that I'd start getting memory failures, etc.  Plus, I think there's a 
leak in the code somewhere.

It was intended to max out a core.

## Can I run this wihtout Go?

Got Docker?  Then yes you can.  Just run:

docker run --rm rharkins/gochan-life [columns] [rows]


