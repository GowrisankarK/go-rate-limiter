# go-rate-limiter
 Rate Limiter developed using Go for Traffic to the Application.

 # Algorithms:
 
 # FixedWindow:
 
    In Fixed Window alogrithm, on the timeframe we allow specific count of requests.

    n request in n time period.

    # Advantages:
    Easy to implement

    # Disadvantages:
    If n requests come in the last second of the current window period and the same n requests come in the first second of the next window period, then there is a chance of a burst.(i.e., Spikes at bounderies)

# SlidingWindow Or RollingWindow:
     In Sliding Window alogrithm, we track the last n duration requests based on which we allow the requests.

    # Advantages:
      Spikes at bounderies won't occurr.
    
    # Disadvantages:
      More Complex than the FixedWindow Algorithm. The reasons are tracking the last n requests count by storing them in memory & cleaning up the older requests.
