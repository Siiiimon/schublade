#include "terminal.h"
#include <cstdio>
#include <unistd.h>

int main(int argc, char** argv) {
    Initialize("/bin/zsh");

    while (1) {
        KernelEvent* e = PollEvent();
        if (e) {
            printf("Event received: %s\n", e->name);
            FreeEvent(e);
        } else {
            usleep(1000 * 100);
        }
    }

    return 0;
}