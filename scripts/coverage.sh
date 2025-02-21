#!/usr/bin/env bash
set -euo pipefail  # Более строгий режим с проверкой неопределенных переменных

# Конфигурация
WORKDIR=".cover"
PROFILE="$WORKDIR/cover.out"
MODE="count"
COVERAGE_HTML="coverage.html"

# Очистка при выходе
cleanup() {
    if [[ "${KEEP_WORKDIR:-}" != "true" ]]; then
        echo "Cleaning up temporary files..."
        rm -rf "$WORKDIR"
    fi
}
trap cleanup EXIT

# Генерация данных покрытия
generate_cover_data() {
    echo "Generating cover data to '$WORKDIR'..."
    rm -rf "$WORKDIR"
    mkdir -p "$WORKDIR"

    local packages
    packages=$(go list ./...)

    for pkg in $packages; do
        local outfile="$WORKDIR/$(echo "$pkg" | tr '/' '-').cover"
        echo "Testing package: $pkg"
        go test -covermode="$MODE" -coverprofile="$outfile" "$pkg" || {
            echo "Error testing package: $pkg"
            return 1
        }
    done

    echo "mode: $MODE" >"$PROFILE"
    grep -h -v "^mode:" "$WORKDIR"/*.cover >>"$PROFILE"
}

# Показ отчета по функциям
show_cover_report_func() {
    echo "Generating functional coverage report..."
    go tool cover -func="$PROFILE"
}

# Генерация HTML отчета
show_cover_report_html() {
    echo "Generating HTML coverage report to '$COVERAGE_HTML'..."
    go tool cover -html="$PROFILE" -o "$COVERAGE_HTML"
}

# Установка прав доступа
set_permissions() {
    if [[ -n "${UID:-}" ]]; then
        echo "Setting UID to $UID..."
        chown -R "$UID" "$COVERAGE_HTML" "$WORKDIR"
    fi

    if [[ -n "${GID:-}" ]]; then
        echo "Setting GID to $GID..."
        chown -R ":$GID" "$COVERAGE_HTML" "$WORKDIR"
    fi
}

main() {
    generate_cover_data
    show_cover_report_func
    show_cover_report_html
    set_permissions
    echo "Coverage analysis complete!"
}

main "$@"