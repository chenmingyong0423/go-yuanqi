SOURCE_COMMIT=.github/pre-commit
TARGET_COMMIT=.git/hooks/pre-commit
SOURCE_PUSH=.github/pre-push
TARGET_PUSH=.git/hooks/pre-push

# copy pre-commit file if not exist.
if [ ! -f $TARGET_COMMIT ]; then
    echo "set git pre-commit hooks..."
    cp $SOURCE_COMMIT $TARGET_COMMIT
fi

# copy pre-push file if not exist.
if [ ! -f $TARGET_PUSH ]; then
    echo "set git pre-push hooks..."
    cp $SOURCE_PUSH $TARGET_PUSH
fi

# add permission to TARGET_PUSH and TARGET_COMMIT file.
test -x $TARGET_PUSH || chmod +x $TARGET_PUSH
test -x $TARGET_COMMIT || chmod +x $TARGET_COMMIT

echo "install golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "install goimports..."
go install golang.org/x/tools/cmd/goimports@latest

echo "go mod tidy"
go mod tidy
