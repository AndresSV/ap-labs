#include <stdio.h>

#define IN   1   /* inside a word */
#define OUT  0   /* outside a word */

/* count lines, words, and characters in input */

void reverse(char array[], int size) {
    for (int i = 0;  i < size/2; i++) {
	    char tmp = array[i];
	    array[i] = array[size - i - 1];
	    array[size - i - 1] = tmp;
    }
    for (int i = 0;  i < size ; i++) {
	    printf("%c" , array[i]);
    }
    printf("\n");
}

int main() {
    int pos = 0;
    char c, 
         word[100];
    while ((c = getchar()) != EOF) {
	    if (c == '\n'){
            reverse(word, pos);
            // for(int i = 0 ; i < pos ; i++){
            //     printf("%c",word[i]);
            // }
            pos = 0;
        } else {
            word[pos] = c;
            pos++;
        }
    }
    return 0;
}