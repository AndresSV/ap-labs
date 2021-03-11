#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <ctype.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>

#define REPORT_FILE "packages_report.txt"

void analizeLog(char *logFile, char *report);
char* getWord(char *word);
char* getDate(char *date);
char* getLine(int l);
char* getHook(char *hook);

struct Package {
  char* name;
  char* installDate;
  char* updateDate;
  int updates;
  char* removalDate;
};

struct Package pkgs[1000];
int installed, upgraded, removed, current;

int main(int argc, char **argv)
{
    if (argc < 3 || argc > 5)
    {
        printf("Invalid number of arguments \n");
    }
    else if((strcmp(argv[1],"-input")==0) && (access(argv[2], R_OK | W_OK )==0) 
        && (strcmp(argv[3],"-report")==0) && (strcmp(argv[4],"")!=0) )
    {
        analizeLog(argv[2], REPORT_FILE);
    }
    else
    {
        printf("Expected:\n");
	    printf("/pacman-analizer -input pacman.txt -report packages_report.txt\n");
    }
    return 0;
}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] \n", logFile);
    int gap = (strcmp(logFile,"pacman2.txt")==0) ? 8:0;
    int r = open(logFile,O_RDONLY);
    char *line;
    int installedPackages   = 0,
        removedPackages     = 0,
        upgradedPackages    = 0,
        currentInstalled    = 0;
    int alpm_scriptletCount = 0,
        alpmCount           = 0,
        pacmanCount         = 0,
        indexPackage        = 0;
    char oldestPackage[256]        = "",
         newestPackage[256]        = "",
         packagesNotUpgraded[1000] = "";
    while((line = getLine(r)) != NULL)
    {
        char *option = getWord(line+26+gap);
        if (strcmp("installed", option) == 0)
        {
            pkgs[installed].name = getWord(line+36+gap);
            pkgs[installed].installDate = getDate(line+1+gap);
            installed++;
            current++;
        }
        else if (strcmp("upgraded", option) == 0)
        {
            char *name =getWord(line+35);
            for (int i = 0; i < installed; i++)
            {
                if (strcmp(pkgs[i].name, name) == 0)
                {
                    pkgs[i].updateDate = getDate(line+1+gap);
                    pkgs[i].updates++;
                    upgraded++;
                }
            }
        }
        else if (strcmp("removed", option) == 0)
        {
            char *name =getWord(line+34+gap);
            for (int i = 0; i < installed; i++)
            {
                if (strcmp(pkgs[i].name, name) == 0)
                {
                    pkgs[i].removalDate = getDate(line+1+gap);
                    removed++;
                    current--;
                }
            }
        }
        char* hook = getHook(line+20+gap);
        if (strcmp(hook, "PACMAN") == 0)
        {
            pacmanCount++;
        }
        else if (strcmp(hook, "ALPM") == 0)
        {
            alpmCount++;
        }
        else if (strcmp(hook, "ALPM-SCRIPTLET") == 0)
        {
            alpm_scriptletCount++;
        }
    }
    for (int i = 0; i < installed; i++)
    {
        if(pkgs[i].removalDate == NULL)
        {
            pkgs[i].removalDate = "-";
        }
        if(pkgs[i].updateDate == NULL)
        {
            pkgs[i].updateDate = "-";
        }
    }
    close(r);

    int fd = open(report, O_WRONLY | O_CREAT, 0600);
    char text[256];
    char* title = "Pacman Packages Report\n----------------------\n";
    write(fd, title, strlen(title));
    sprintf(text, "%s%d\n", "- Installed packages : ", installed);
    write(fd, text, strlen(text));
    sprintf(text, "%s%d\n", "- Removed packages : ", removed);
    write(fd, text, strlen(text));
    sprintf(text, "%s%d\n", "- Upgraded packages : ", upgraded);
    write(fd, text, strlen(text));
    sprintf(text, "%s%d\n", "- Current installed : ", current);
    write(fd, text, strlen(text));
    snprintf(text,100,"-------------\nGeneral Stats\n-------------\n");
    write(fd, text, strlen(text));
    snprintf(text,100,"- Oldest package               : %s \n",oldestPackage);
    write(fd, text, strlen(text));
    snprintf(text,100,"- Package with no upgrades     : %s \n", packagesNotUpgraded);
    write(fd, text, strlen(text));
    snprintf(text,100,"- [ALPM-SCRIPTLET] type count  : %d \n",alpm_scriptletCount);
    write(fd, text, strlen(text));
    snprintf(text,100,"- [ALPM] count                 : %d \n",alpmCount);
    write(fd, text, strlen(text));
    snprintf(text,100,"- [PACMAN] count               : %d \n",pacmanCount);
    write(fd, text, strlen(text));

    char* secondTitle = "----------------\nList of packages\n----------------\n";
    write(fd, secondTitle, strlen(secondTitle));

    for (int i = 0; i < installed; i++)
    {
        sprintf(text, "%s%s\n", "- Package Name : ", pkgs[i].name);
        write(fd,text,strlen(text));
        sprintf(text, "%s%s\n", "     - Install date : ", pkgs[i].installDate);
        write(fd,text,strlen(text));
        sprintf(text, "%s%s\n", "     - Last update date : ", pkgs[i].updateDate);
        write(fd,text,strlen(text));
        sprintf(text, "%s%d\n", "     - How many updates : ", pkgs[i].updates);
        write(fd,text,strlen(text));
        sprintf(text, "%s%s\n", "     - Removal date : ", pkgs[i].removalDate);
        write(fd,text,strlen(text));
    }
    close(fd);
    printf("Report is generated at: [%s]\n", report);
}

char* getWord(char *word){
    int size;
    for(int i=0; i<strlen(word); i++)
    {
        if(word[i] == ' ' || word[i] == '\n')
        {
            size = i;
            break;
        }
    }
    char *wordRes = calloc(size+1,1);
    for(int i=0; i<size; i++)
    {
        wordRes[i] = word[i];
    }
    wordRes[size] = '\0';
    return wordRes;
}

char* getDate(char *date){
    int size;
    for(int i=0;i<strlen(date);i++)
    {
        if(date[i] == ']')
        {
            size = i;
            break;
        }
    }
    char *dateRes = calloc(size+1,1);
    for(int i=0; i<size; i++)
    {
        dateRes[i] = date[i];
    }
    dateRes[size] = '\0';
    return dateRes;
}

char* getLine(int l)
{
    int size = 256,
        total = 0,
        n;
    char *buffer;
    buffer = calloc(size, 1);
    while((n = read(l, buffer+total, size)) > 0)
    {
        for(int i=total; i<total+n; i++)
        {
            if(buffer[i] == '\n')
            {
                lseek(l, i-(total+n)+1, SEEK_CUR);
                buffer[i] = '\0';
                return buffer;
            }
        }
        total+=n;
        buffer = realloc(buffer, total+size);
    }
    if( total > 0 )
    {
        buffer[total] = '\0';
        return buffer;
    }
    return NULL;
}

char* getHook(char *hook){
    int size;
    for(int i=0;i<strlen(hook);i++)
    {
        if(hook[i] == ']')
        {
            size = i;
            break;
        }
    }
    char *hookRes = calloc(size+1,1);
    for(int i=0; i<size; i++)
    {
        hookRes[i] = hook[i];
    }
    hookRes[size] = '\0';
    return hookRes;
}