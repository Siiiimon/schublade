#include "terminal.h"
#include <cstdio>

int main(int argc, char** argv) {
    Initialize("/bin/zsh");

    while (1) {
        KernelEvent* e = PollEvent();
        if (e) {
            printf("Event received: %d - %s: '%s'\n", e->kind, e->name, (char*) e->data);
            FreeEvent(e);
        }
    }

    return 0;
}