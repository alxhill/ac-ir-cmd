#include "irslinger.h"

int main(int argc, char *argv[]) {
    if (argc < 2) {
        return -1;
    }

    char *cmd = argv[1];

    return irSling(
        17,
        38000
        0.5,
        9000
        4500,
        562,
        562,
        1688,
        562,
        1,
        cmd);
}