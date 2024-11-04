cmake_minimum_required(VERSION 3.29)
project(schublade)

set(CMAKE_CXX_STANDARD 23)

add_executable(schublade main.cpp)

if(WIN32)
    target_link_libraries(${PROJECT_NAME} PRIVATE Dwmapi.lib)
endif()