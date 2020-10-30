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

## Test results (Nov 2020)

* A simulated ACU datastream compressed by a factor of ~1.6.
* Limited support for encodings:
    - Neither `pyarrow` nor `fastparquet` can read `DELTA_BINARY_PACKED`.
    - Go can't write `BYTE_STREAM_SPLIT`.
* `fastparquet` can't handle zstd compression.

