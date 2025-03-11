#include <iostream>
#include <fstream>
#include <curl/curl.h>

size_t WriteCallback(void* contents, size_t size, size_t nmemb, std::string* output) {
    output->append((char*)contents, size * nmemb);
    return size * nmemb;
}

int main() {
    CURL* curl;
    CURLcode res;
    std::string readBuffer;

    curl = curl_easy_init();
    if(curl) {
        curl_easy_setopt(curl, CURLOPT_URL, "http://localhost:8080/generate");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, "{\"language\": \"Python\"}");
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);

        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);

        std::ofstream outFile("python-course.txt");
        outFile << readBuffer;
        outFile.close();

        std::cout << "✅ Course saved to python-course.txt" << std::endl;
    }
    return 0;
}
