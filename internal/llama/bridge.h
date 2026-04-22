#ifndef GO_LLAMA_CPP_LAB_BRIDGE_H
#define GO_LLAMA_CPP_LAB_BRIDGE_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

int bridge_llama_backend_init(char * err_buf, size_t err_buf_size);
void bridge_llama_backend_free(void);
int bridge_llama_max_devices(size_t * out, char * err_buf, size_t err_buf_size);
int bridge_llama_print_system_info(char * buf, size_t buf_size, char * err_buf, size_t err_buf_size);

#ifdef __cplusplus
}
#endif

#endif
