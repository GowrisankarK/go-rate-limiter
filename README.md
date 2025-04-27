# go-rate-limiter

A Rate Limiter developed in Go to control traffic to the application.

## Algorithms

### Fixed Window

In the **Fixed Window** algorithm, within a fixed time frame, we allow a specific number of requests.

**Example:**  
Allow `n` requests every `n` milliseconds.

#### Advantages
- Easy to implement.

#### Disadvantages
- If `n` requests come at the **end** of the current window and another `n` requests at the **start** of the next window, it can cause a **burst** (i.e., spikes at boundaries).

---

### Sliding Window Log (Rolling Window Log)

In the **Sliding Window Log** algorithm, we track the timestamps of individual requests over the last `n` duration to decide if new requests should be allowed.

#### Advantages
- No spikes at boundaries.

#### Disadvantages
- More complex than the Fixed Window algorithm.
- Higher **space complexity**, as it stores **each request's timestamp**.
- Requires periodic cleanup of old request timestamps.

---

### Sliding Window Counter (Rolling Window Counter)

The **Sliding Window Counter** is similar to the Sliding Window Log, but instead of storing every request's timestamp, we group requests by their timestamps (e.g., by second or millisecond) and count occurrences, using a key-value pair data structure.

#### Advantages
- Reduces space complexity compared to Sliding Window Log.
- No spikes at boundaries.

#### Disadvantages
- Still more complex than the Fixed Window algorithm.
- Requires logic to track counts and clean up old timestamps.

---

### Token Bucket

In the **Token Bucket** algorithm, tokens are refilled into a bucket at a regular interval.  
If a token is available when a request comes, the request is allowed and a token is consumed.

#### Advantages
- Prevents spikes at boundaries.
- Can handle **bursty** traffic while maintaining an average rate over time.

#### Disadvantages
- More complex than the Fixed Window algorithm.
- Needs careful implementation for token refill timing and accurate synchronization.

---

### Leaky Bucket

In the **Leaky Bucket** algorithm, incoming requests are placed into a **fixed-size queue**.  
Requests are processed at a fixed rate. If the queue is full, new requests are dropped.

#### Advantages
- Smooths traffic flow and prevents spikes at boundaries.
- Simple to understand and manage.

#### Disadvantages
- More complex than the Fixed Window algorithm.
- Requests can be dropped if the **incoming rate exceeds the processing rate** for a long time.

---

## Algorithms Comparison

| Algorithm                  | Memory Usage             | Complexity Level          | Supports Bursts?  | Pros                                         | Cons                                            |
|-----------------------------|---------------------------|----------------------------|-------------------|----------------------------------------------|-------------------------------------------------|
| **Fixed Window**            | Very Low                  | Very Easy                  | ❌ No              | Easy to implement                            | Bursts at window boundaries                    |
| **Sliding Window Log**      | High (stores timestamps)   | Moderate                   | ❌ No              | Smooth traffic, prevents boundary spikes    | High memory usage, cleanup needed              |
| **Sliding Window Counter**  | Medium (buckets timestamps)| Moderate                   | ❌ No              | Smoother traffic with lower memory          | Approximation errors possible                  |
| **Token Bucket**            | Low (token counter only)   | High                       | ✅ Yes             | Handles sudden bursts, adjustable traffic    | Complex refill logic, time sync needed          |
| **Leaky Bucket**            | Low (queue based)          | Moderate                   | ❌ No              | Smooths traffic flow, simple to reason about | Drops requests if incoming rate too high       |

---