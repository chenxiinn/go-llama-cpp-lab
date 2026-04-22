//go:build llama
// +build llama

#include "bridge.h"

#include "llama.h"

#include <stdio.h>
#include <string.h>

static int bridge_set_error(char * err_buf, size_t err_buf_size, const char * msg) {
    if (err_buf != NULL && err_buf_size > 0) {
        snprintf(err_buf, err_buf_size, "%s", msg);
    }
    return -1;
}

int bridge_llama_backend_init(char * err_buf, size_t err_buf_size) {
    (void) err_buf;
    (void) err_buf_size;

    llama_backend_init();
    return 0;
}

void bridge_llama_backend_free(void) {
    llama_backend_free();
}

int bridge_llama_max_devices(size_t * out, char * err_buf, size_t err_buf_size) {
    if (out == NULL) {
        return bridge_set_error(err_buf, err_buf_size, "max devices output pointer is null");
    }

    *out = llama_max_devices();
    return 0;
}

int bridge_llama_print_system_info(char * buf, size_t buf_size, char * err_buf, size_t err_buf_size) {
    const char * info;

    if (buf == NULL || buf_size == 0) {
        return bridge_set_error(err_buf, err_buf_size, "system info buffer is empty");
    }

    info = llama_print_system_info();
    if (info == NULL) {
        return bridge_set_error(err_buf, err_buf_size, "llama_print_system_info returned null");
    }

    snprintf(buf, buf_size, "%s", info);
    return 0;
}
