#include <stdio.h>
#include <string.h>
#include <fstream>
#include <iostream>
 
#include <curl/curl.h>

using namespace std;


char getImage(char *&data, char *filename, unsigned int &len){
  std::ifstream is (filename, std::ifstream::binary);
  if (is) {
    // get len of file:
    is.seekg (0, is.end);
    len = is.tellg();
    is.seekg (0, is.beg);

    data = new char [len];

    std::cout << "Reading " << len << " characters... ";
    // read data as a block:
    is.read (data,len);

    if (is)
      std::cout << "all characters read successfully.";
    else{
      std::cout << "error: only " << is.gcount() << " could be read";
      return 0;
    }
    is.close();
    return 1;
  }
  return 0;
}

 
int main(int argc, char *argv[]){

  curl_mime *form = NULL;
  curl_mimepart *field = NULL;
  struct curl_slist *headerlist = NULL;
  static const char buf[] = "Expect:";
  CURL *curl;
  CURLcode res;
  char *filename;
  char *filenameToSend;
  char *authtoken;
  char *url;
  int port;
  char *data;
  char ext[20];
  unsigned int len;

  if(argc != 6){
    fprintf(stderr, "Error\nUsage: program <file> <filename to send> <authtoken> <url> <port>\n");
    exit(0);
  }
  
  filename = argv[1];
  filenameToSend = argv[2];
  authtoken = argv[3];
  url = argv[4];
  port = atoi(argv[5]);
  printf("%s\n", authtoken);

  curl_global_init(CURL_GLOBAL_ALL);
 
  curl = curl_easy_init();
  if(curl) {
    /* Create the form */ 
    form = curl_mime_init(curl);
 
    char ok = getImage(data, filename, len);

    if(!ok){
      printf("Error getting the image\n");
      exit(0);
    }

    /* Fill in the file upload field */ 
    field = curl_mime_addpart(form);
    curl_mime_name(field, "sendfile");
    curl_mime_data(field, data, len);
    curl_mime_filename(field, filenameToSend);
    curl_mime_type(field, ext);
    
    field = curl_mime_addpart(form);
    curl_mime_name(field, "authtoken");
    curl_mime_data(field, authtoken, CURL_ZERO_TERMINATED);
 
    headerlist = curl_slist_append(headerlist, buf);

    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0);
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_PORT, port);

    curl_easy_setopt(curl, CURLOPT_MIMEPOST, form);
    // curl_easy_setopt(curl, CURLOPT_POSTFIELDS, "authtoken=asdf");
 
    /* Perform the request, res will get the return code */ 
    res = curl_easy_perform(curl);
    /* Check for errors */ 
    if(res != CURLE_OK)
      fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));
 
    /* always cleanup */ 
    curl_easy_cleanup(curl);
 
    /* then cleanup the form */ 
    curl_mime_free(form);
    /* free slist */ 
    curl_slist_free_all(headerlist);
    delete[] data;
  }
  return 0;
}