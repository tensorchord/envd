set -eu

echo 'zip all file started with envd'
echo

ls | while read last; do
    if [[ "$last" == envd* ]]
    then
        gzip $last
        echo 'zip '$last
    fi
done
