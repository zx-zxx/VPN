# 1. Introdunction and Background

state the objective and background



# 2. Methodology

Experimental Setup

Tools

Procedure



# 3. Results

including the filled Table 1



| Metric                     | Phase         | TLS 1.3 (Baseline) | PQC TLS | PQC TLS（MLKEM768） | Units | Notes                       |
| -------------------------- | ------------- | ------------------ | ------- | ------------------- | ----- | --------------------------- |
| Handshake Duration         | Handshake     | 1.6                | 2.2     | 1.3                 | ms    | Wireshark, 5 trials         |
| Throughput                 | Data Transfer | \（没测成功）      | \       | \                   | Mbps  | iperf3 average, 5 trials    |
| Client CPU Usage           | Handshake     | 43.0               | 111.3   | 122.9               | %     | docker stats peak, 5 trials |
| Client CPU Usage           | Data Transfer | 3.9                | 1.64    | 1.26                | %     | pidstat avg, 5 trials       |
| Server CPU Usage           | Handshake     | 0.22               | 0.25    | 0.21                | %     | docker stats peak, 5 trials |
| Server CPU Usage           | Data Transfer | 0.79               | 0.63    | 0.22                | %     | pidstat avg, 5 trials       |
| Client Memory Usage (Heap) | Handshake     | 17.82              | 31.4    | 21.32               | MB    | docker stats peak, 5 trials |
| Client Memory Usage (Heap) | Data Transfer | 18.2               | 19.31   | 25.04               | MB    | docker stats avg, 5 trials  |
| Server Memory Usage (Heap) | Handshake     | 2.71               | 2.85    | 3.46                | MB    | docker stats peak, 5 trials |
| Server Memory Usage (Heap) | Data Transfer | 2.69               | 3.0     | 3.68                | MB    | docker stats avg, 5 trials  |

# 4. Analysis

Comparison

Interpretation of Data and pprof profiles

exploring the security vs. performance trade-offs



# 5. Conclusion

 summarizing your findings

