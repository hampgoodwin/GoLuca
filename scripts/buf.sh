if ! command -v brew &> /dev/null
then
    echo "brew could not be found"
    exit 1
fi

if ! command -v buf &> /dev/null
then
    echo "attemping global install for buf cli"
    brew install bufbuild/buf/buf
    echo "brew could not be found"
    exit 1
fi