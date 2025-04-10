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


# TokenBucket:
    In Token Bucket alogrithm, we refill the tokens count in bucket on specific period. If the token is available in a bucket, the request is allowed and token is reduced by one.

    # Advantages:
      Spikes at bounderies won't occurr.
      Can be adjusted to handle varying traffic patterns
    
    # Disadvantages:
      More Complex than the FixedWindow Algorithm.

# LeakyBucket:
    In Leaky Bucket alogrithm, we have a queue with max size to hold the request. on specific period, the request is consumed. If the queue is full, then request won't be accepted.

    # Advantages:
      Spikes at bounderies won't occurr.
      Simple to understand and manage.
    
    # Disadvantages:
      More Complex than the FixedWindow Algorithm.
      Can lead to dropped requests if the incoming rate consistently exceeds the processing rate.