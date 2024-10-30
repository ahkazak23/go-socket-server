#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock2.h>  // Use winsock2.h for Windows
#include <ws2tcpip.h>  // For inet_pton and other TCP/IP functions

#define PORT 8080
#define BUFFER_SIZE 1024

int main() {
    WSADATA wsaData;
    int sock = 0;
    struct sockaddr_in server_addr;
    char buffer[BUFFER_SIZE];

    // Initialize Winsock
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        printf("Failed to initialize Winsock. Error Code: %d\n", WSAGetLastError());
        return -1;
    }

    // Create socket
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == INVALID_SOCKET) {
        printf("Socket creation error. Error Code: %d\n", WSAGetLastError());
        WSACleanup();
        return -1;
    }

    // Set up server address struct
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);

    // Convert IPv4 address from text to binary
    if (inet_pton(AF_INET, "127.0.0.1", &server_addr.sin_addr) <= 0) {
        printf("Invalid address/ Address not supported\n");
        closesocket(sock);
        WSACleanup();
        return -1;
    }

    // Connect to the server
    if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
        printf("Connection Failed. Error Code: %d\n", WSAGetLastError());
        closesocket(sock);
        WSACleanup();
        return -1;
    }

    // Read welcome message from server
    memset(buffer, 0, BUFFER_SIZE);
    int valread = recv(sock, buffer, BUFFER_SIZE, 0);
    if (valread > 0) {
        printf("Server: %s\n", buffer);
    }

    // Interactive loop to send commands and read responses
    while (1) {
        // Get user input for the next command
        fgets(buffer, BUFFER_SIZE, stdin);

        // Ensure the input ends with a newline
        size_t len = strlen(buffer);
        if (len > 0 && buffer[len - 1] != '\n') {
            buffer[len] = '\n';  // Add newline if not present
            buffer[len + 1] = '\0';  // Null-terminate the string
        }

        // Send the user input to the server
        int bytes_sent = send(sock, buffer, strlen(buffer), 0);
        if (bytes_sent < 0) {
            printf("Error sending data.\n");
        }

        // Exit the loop if the user typed 'exit'
        if (strcmp(buffer, "exit\n") == 0) {
            break;
        }

        // Read and display the response from the server
        memset(buffer, 0, BUFFER_SIZE);
        valread = recv(sock, buffer, BUFFER_SIZE, 0);
        if (valread > 0) {
            printf("Server: %s\n", buffer);
        } else {
            printf("Server closed the connection.\n");
            break;
        }
    }

    // Close the socket and clean up Winsock
    closesocket(sock);
    WSACleanup();
    return 0;
}
