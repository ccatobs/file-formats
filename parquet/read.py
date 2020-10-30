#!/usr/bin/env python3

import sys
import time

import numpy as np
import pandas as pd

import pyarrow.parquet as pq
#import fastparquet as pq

filename = sys.argv[1]

start = time.time()
pf = pq.ParquetFile(filename)
if 'pyarrow' in sys.modules:
    df = pf.read().to_pandas()
else:
    df = pf.to_pandas()
elapsed_secs = time.time() - start

nrows = pf.metadata.num_rows
row_bytes = int(pf.metadata.metadata[b'row_bytes'])

print(df)
print()
print(f"read speed = {1e-6*nrows*row_bytes/elapsed_secs} MB/s")

