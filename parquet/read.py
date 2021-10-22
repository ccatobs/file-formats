#!/usr/bin/env python3

import sys
import time

import numpy as np
import pandas as pd

which = sys.argv[1]
if which == 'pyarrow':
    import pyarrow.parquet as pq
else:
    import fastparquet as pq

filename = sys.argv[2]

start = time.time()
pf = pq.ParquetFile(filename)

if which == 'pyarrow':
    df = pf.read().to_pandas()
    nrows = pf.metadata.num_rows
    row_bytes = int(pf.metadata.metadata[b'row_bytes'])
else:
    df = pf.to_pandas()
    nrows = pf.info['rows']
    row_bytes = int(pf.key_value_metadata['row_bytes'])

elapsed_secs = time.time() - start

print(df)
print()
print(f"read speed = {1e-6*nrows*row_bytes/elapsed_secs} MB/s")

