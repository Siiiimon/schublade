cmake_minimum_required(VERSION 3.29)
project(schublade)
set(CMAKE_CXX_STANDARD 23)

add_executable(schublade main.cpp)

# build go kernel

if (WIN32)
    set(GO_LIBRARY "${CMAKE_CURRENT_BINARY_DIR}/terminal.dll")
elseif(APPLE)
    set(GO_LIBRARY "${CMAKE_CURRENT_BINARY_DIR}/terminal.dylib")
endif()
set(GO_SOURCE_FILE
        "${CMAKE_CURRENT_SOURCE_DIR}/kernel/kernel.go"
        "${CMAKE_CURRENT_SOURCE_DIR}/kernel/event.go"
        "${CMAKE_CURRENT_SOURCE_DIR}/kernel/logger.go"
        "${CMAKE_CURRENT_SOURCE_DIR}/kernel/parser.go"
        "${CMAKE_CURRENT_SOURCE_DIR}/kernel/shell.go"
)

add_custom_command(
        OUTPUT ${GO_LIBRARY}
        COMMAND go build -o ${GO_LIBRARY} -buildmode=c-shared ${GO_SOURCE_FILE}
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/kernel
        COMMENT "Building Kernel"
        DEPENDS ${GO_SOURCE_FILE}
)
add_custom_target(BuildGoLibrary ALL DEPENDS ${GO_LIBRARY})
add_dependencies(${PROJECT_NAME} BuildGoLibrary)


# build gui
include_directories("${CMAKE_CURRENT_BINARY_DIR}")

if(WIN32)
    include_directories("${CMAKE_CURRENT_BINARY_DIR}")

    target_link_libraries(${PROJECT_NAME} PRIVATE
            d2d1 dwrite
            "${CMAKE_CURRENT_BINARY_DIR}/terminal.dll"
    )
elseif(APPLE)
    target_link_libraries(${PROJECT_NAME} PRIVATE
            "${CMAKE_CURRENT_BINARY_DIR}/terminal.dylib"
    )
endif()