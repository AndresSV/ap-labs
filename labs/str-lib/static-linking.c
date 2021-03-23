#include <stdio.h>
#include <stdlib.h>

int mystrlen(char *);
char *mystradd(char *,char *);
int mystrfind(char *,char *);

int main(int argc, char *argv[]) {
    char *string = (char *)malloc(sizeof(char *)),
         *string2 = (char *)malloc(sizeof(char *)),
         *r = (char *)malloc(sizeof(char *));
    string = argv[2];
    string2 = argv[3];
    if (mystrfind(argv[1],"-add")==0 && argc == 4)
    {
        printf("Initial Lenght\t: %d\n", mystrlen(string));
        mystradd(string,string2);
        printf("New String\t: %s\n", string);
        printf("New lenght\t: %d\n", mystrlen(string));
        return 0;
    }
    else if(mystrfind(argv[1],"-find")==0 && argc == 4)
    {
        printf("[%s] string was found at [%d] position\n", string2,mystrfind(string,string2));
    }
    else
    {
        printf("Expected:\n");
        printf("./main.o -add \"[string]\" \"[string]\"\n");
        printf("- or -");
        printf("./main.o -find \"[string]\" \"[string]\"\n");
    }
    return 0;
}
