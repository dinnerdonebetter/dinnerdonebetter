#!/bin/bash

. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/lib/ensure_binaries.sh

THIS="github.com/dinnerdonebetter/backend"
WIRE_TARGETS=(
    "internal/build/services/api"
    "internal/build/services/admin_webapp"
    "internal/build/jobs/db_cleaner"
    "internal/build/jobs/search_data_index_scheduler"
    "internal/build/jobs/meal_plan_finalizer"
    "internal/build/jobs/meal_plan_grocery_list_initializer"
    "internal/build/jobs/meal_plan_task_creator"
)

rewire() {
    ensure_wire_installed
    for tgt in "${WIRE_TARGETS[@]}"; do
        target_path="${THIS}/${tgt}"
        echo "Rewiring ${target_path}"
        rm -f "${target_path}/wire_gen.go"
        wire gen "${target_path}"
    done
    # TODO: format
}

rewire