A simple and fast CLI for checking CSV files for errors.


## Usage

```text
checkcsv [ -options ] < <file.csv>
```


## Options

```text
-c             Collect all csv errors and output the list at the end
-s             Print summary info (check status, total columns, total rows)
-sj            Print summary infoas a JSON string
-q             Silently terminate with exit(1) upon the first error encountered in the CSV
-d             Fields separator (default: comma)
-h, --help     Help
-v, --version  Version
```


## Examples
```shell
# List all errors
checkcsv -c < file.csv

# Check gzipped csv file
gunzip -c "file.csv.gz" | checkcsv && echo "ok" || echo "(!) error"

# Get check summary as a JSON
checkcsv -sj < file.csv
```




