#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
    if (argc == 2) {
        float fahrenheit = atoi(argv[1]);
        float celcsius   = (fahrenheit - 32) * 5 / 9 ;
        printf("Fahrenheit: %3.0f, Celcius: %6.1f\n", fahrenheit, celcsius);
    }
    else if (argc == 4)
    {
        float fahrenheit;
        int LOWER = atoi(argv[1]),
            UPPER = atoi(argv[2]), 
            STEP  = atoi(argv[3]);
        for (fahrenheit = LOWER; fahrenheit <= UPPER; fahrenheit = fahrenheit + STEP)
        {
            float celcsius   = (fahrenheit - 32) * 5 / 9 ;
            printf("Fahrenheit: %3.0f, Celcius: %6.1f\n", fahrenheit, celcsius);
        }
    }
    return 0;
}