# Scaler

Scaler measures how Kask's total allocations scale with content directory sizes. It repeatedly calls the builder on a content directory which are generated during the session at sizes linarly grow in bytes. After each run of builder Scalers notes down the `runtime.MemStat.Sys` value as the total allocations.

The test site content is based on the Kask docs site. Scaler only grows the sample site with the very same docs site content; so the test site contents reflects the average real world case at each size.

Scaler grows the sample site size at each iteration drammatically in order to sample the Kask's behavior in greater range and avoid confusing by Go allocating more a program may need.

At last, Scaler performs a heuristic and returns with the status code. Success means Kask allocates memory from the system grows sublinearly with the content directory size.

Scaler designed to be used in the CI pipeline; inside a separate workflow than the push/PR triggered main workflow to perform checks occasionally and warn when a regression happens without constant CI wait and noise.
