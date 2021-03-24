#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int numcmp(char *s1, char *s2) {
    return atof(s1)-atof(s2);
}

void merge(char *lineptr[], int l, int m, int r, int (*comp)(void *,void *))
{
	int i = 0,
        j = 0,
        k = 0,
        s1 = m-l+1,
        s2 = r-m;

	void *left[s1],
         *right[s2];

	for (i=0; i<s1; i++)
    {
		left[i] = lineptr[l+i];
	}
	for (i=0; i<s2; i++)
    {
		right[i] = lineptr[m+i+1];
	}
	i = j = 0;
    k = l;
    while(i<s1 && j<s2){
    	if(comp(left[i],right[j]) < 0)
        {
            lineptr[k++] = left[i++];
		}
        else
        {
            lineptr[k++] = right[j++];
		}
	}
	if (i == s1)
    {
	    while(j < s2)
        {
			lineptr[k++] = right[j++];
		}
	}
	else{
		while(i<s1){
			lineptr[k++] = left[i++];
		}
	}
}

void mergeSort(char *lineptr[], int start, int end, int (*comp)(void *,void *))
{
	int middle = start+(end-start)/2;
	if(end>start){
		mergeSort(lineptr, start, middle, comp);
		mergeSort(lineptr, middle+1, end, comp);
		merge(lineptr, start, middle, end, comp);
	}
}

void writeOut(char *file, int flag)
{
    int arraySize = 100;
    char **arrayList = malloc(sizeof(char *) * arraySize);

    printf("%s\n",file);
    FILE *fp = fopen(file, "r");
    if (!fp){
        printf("Error, Failed to open the file! \n");
    }

    char *line_buf = NULL;
    size_t line_buf_size = 0;
    int line_count = 0;
    ssize_t line_size;

    line_size = getline(&line_buf, &line_buf_size, fp);
    while (line_size >= 0)
    {
        arrayList[line_count] = malloc(sizeof(char) * line_buf_size);
        sprintf(arrayList[line_count], "%s", line_buf);
        arrayList[line_count][strlen(line_buf) - 1] = 0;
        line_count++;
        line_size = getline(&line_buf, &line_buf_size, fp);
    }

    free(line_buf);
    fclose(fp);

    if(flag == 0){
        mergeSort(arrayList, 0, arraySize - 1, (int (*)(void *, void *)) strcmp);
    }
    else if(flag ==1){
        mergeSort(arrayList, 0, arraySize - 1, (int (*)(void *, void *)) numcmp);
    }

    FILE *fpOut;
    fpOut = fopen((flag?"sorted_numbers.txt":"sorted_strings.txt"),"w");
    if (!fpOut){
        printf("Error, Failed to open the file! \n");
    }
    
    for(int i=0; i<100; i++)
    {
        fprintf(fpOut, "%s\n",arrayList[i]);
    }
    fclose(fpOut);
}

int main(int argc, char **argv)
{
    if (argc == 2 && strstr(argv[1],".txt")!=NULL)
    {
        writeOut(argv[1],0);
        printf("Results file can be found at ./sorted_strings.txt\n");
        return 0;
    }
    else if(argc == 3 && strcmp(argv[1], "-n")==0 && strstr(argv[2],".txt")!=NULL)
    {
        writeOut(argv[2],1);
        printf("Results file can be found at ./sorted_numbers.txt\n");
        return 0;
    }
    else 
    {
        printf("Expected:\n");
        printf("$ ./generic_merge_sort.o -n numbers.txt\n");
        printf("$ ./generic_merge_sort.o strings.txt\n");
        return 0;
    }
}