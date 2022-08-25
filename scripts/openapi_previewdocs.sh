if ! command -v npm &> /dev/null
then
    echo "npm could not be found"
    exit 1
fi

if ! command -v openapi &> /dev/null
then
    echo "attemping global install for openapi cli"
    npm i -g @redocly/openapi-cli@latest
    echo "<the_command> could not be found"
    exit 1
fi

openapi preview-docs http/spec/openapi.yml
