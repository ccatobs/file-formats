# Parquet file format

## About

Pros:

* Stores data in columns, improving compression.
* Handles arbitrarily nested data.
* Emerging standard for the Apache "big data" ecosystem.
* Good python/pandas integration.

Cons:

* [Website](https://parquet.apache.org) out of date.
* A little unclear who's in charge. Has the [Arrow](https://arrow.apache.org) project has taken over?
* Arrow development is all in one giant monorepo.
* [Compatibility issues between Spark and Arrow](https://issues.apache.org/jira/browse/ARROW-6057)?

## Test results (Nov 2020)

* A simulated ACU datastream compressed by a factor of ~1.6 with zstd.
* Limited support for encodings:
    - Neither `pyarrow` nor `fastparquet` can read `DELTA_BINARY_PACKED`.
    - Go can't write `BYTE_STREAM_SPLIT`.
* `fastparquet` can't handle zstd compression.
* On Linux, running in docker and writing to a bind mount is ~5% slower versus running natively.
  On macOS, the slowdown is much larger (~80%).

## Test results (Oct 2021)

* `fastparquet` can now handle zstd compression.
* `pyarrow` is much faster at reading than `fastparquet`.

