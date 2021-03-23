#include <stdlib.h>

int mystrlen(char *str){
    int count = 0;
    while (str[count] != '\0')
    {
        count++;
    }
    return count;
}

char *mystradd(char *origin, char *addition){
    int count       = mystrlen(origin),
        additionlen = mystrlen(addition);
    for(int i = 0; i < additionlen; i++)
    {
        *(origin+count++) = *(addition+i);
    }
    *(origin+count) = '\0';
    return origin;
}

int mystrfind(char *origin, char *substr){
    int count       = 0,
        limit       = mystrlen(substr),
        originlen   = mystrlen(origin);
    for(int i = 0; i < originlen; i++)
    {
        if(*(origin+i) == *(substr+count))
        {
            count++;
        }
        else 
        {
            i = i - count + 1;
            count = 0;
        }
        if(count==limit){
            return i - count + 1;
        }
    }
    return -1;
}