#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* month_day function's prototype*/
void month_day(int year, int yearday, int *pmonth, int *pday);
/* month_name:  return name of n-th month */
char *month_name(int n);

int main(int argc, char *argv[]) {
    
    if(argc != 3)
    {
        printf("Expected.\n");
        printf("./month_day <year> <yearday>\n");
    }
    else{
        int year    = atoi(argv[1]),
            yearday = atoi(argv[2]),
            *pmonth = malloc(sizeof(int)),
            *pday = malloc(sizeof(int));
        month_day(year, yearday, pmonth, pday);
        char *month = month_name(*pmonth);
        if (strcmp(month, "Illegal month")!=0)
        {
            printf("%s %d, %d \n" ,month_name(*pmonth), *pday, year);
        }
    }
    return 0;

}

void month_day(int year, int yearday, int *pmonth, int *pday)
{
    static char daytab[2][13] = {
                                    {0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31},
                                    {0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
                                };    
    int leap = ( ( (year%4 == 0) && (year%100 != 0) ) || (year%400 == 0) ) ;
    if (year < 1 || yearday <= 0 || yearday > ((leap == 0)?365:366) )
    {
        printf("Incorrect day.\n");
        return;
    }
    int daySum = 0, i = 0, prevDaySum = 0;
    while (daySum < yearday)
    {
        prevDaySum = daySum;
        daySum += daytab[leap][i];
        ++i;
    }
    *pmonth = --i;
    *pday   = yearday - prevDaySum;
}

char *month_name(int n) {
    static char *name[] = {
                            "Illegal month",
                                "January", "February", "March",
                                "April", "May", "June",
                                "July", "August", "September",
                                "October", "November", "December"
                            };
   return (n < 1 || n > 12) ? name[0] : name[n];
}