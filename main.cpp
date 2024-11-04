#include <iostream>
#include <windows.h>

#define OPACITY 90
#define WIDTH_PERCENTAGE 0.8
#define HEIGHT_PERCENTAGE 0.5
#define DROPDOWN_DISTANCE_IN_PIXEL_PER_TICK 20
#define DROPDOWN_DURATION_IN_MILLISECONDS 200

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
    static bool isAnimating = false;

    switch (uMsg)
    {
        case WM_KEYDOWN: {
            if (wParam == VK_ESCAPE) {
                PostMessage(hwnd, WM_CLOSE, 0, 0);
            }
            return 0;
        }
        case WM_TIMER: {
            if (isAnimating) {
                RECT rect;
                GetWindowRect(hwnd, &rect);
                int currentY = rect.top;

                if (currentY < 0) {
                    SetWindowPos(hwnd, NULL, rect.left, currentY + DROPDOWN_DISTANCE_IN_PIXEL_PER_TICK, 0, 0, SWP_NOSIZE | SWP_NOZORDER);
                } else {
                    KillTimer(hwnd, 1);
                    isAnimating = false;
                }
            }
            return 0;
        }
        case WM_SHOWWINDOW: {
            if (wParam) {
                int screenWidth = GetSystemMetrics(SM_CXSCREEN);
                int screenHeight = GetSystemMetrics(SM_CYSCREEN);

                int windowWidth = static_cast<int>(screenWidth * 0.8);
                int windowHeight = static_cast<int>(screenHeight * 0.5);
                int windowX = (screenWidth - windowWidth) / 2;

                SetWindowPos(hwnd, NULL, windowX, -windowHeight, windowWidth, windowHeight, SWP_NOZORDER);

                static int ticks = windowHeight / DROPDOWN_DISTANCE_IN_PIXEL_PER_TICK;

                isAnimating = true;
                SetTimer(hwnd, 1, DROPDOWN_DURATION_IN_MILLISECONDS / ticks, NULL);
            }
            return 0;
        }
        case WM_PAINT: {
            PAINTSTRUCT ps;
            HDC hdc = BeginPaint(hwnd, &ps);

            HBRUSH hBrush = CreateSolidBrush(RGB(34, 37, 48));
            FillRect(hdc, &ps.rcPaint, hBrush);
            DeleteObject(hBrush);

            EndPaint(hwnd, &ps);
            return 0;
        }
        case WM_DESTROY: {
            PostQuitMessage(0);
            return 0;
        }
        default: {
            return DefWindowProc(hwnd, uMsg, wParam, lParam);
        }
    }
}

int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance,
    PSTR lpCmdLine, int nCmdShow)
{
    LPCSTR CLASS_NAME  = "SchubladeWindowClass";

    WNDCLASS wc = {};
    wc.lpfnWndProc = WindowProc;
    wc.hInstance = hInstance;
    wc.lpszClassName = CLASS_NAME;

    RegisterClass(&wc);

    HWND hwnd = CreateWindowEx(
        WS_EX_TOPMOST | WS_EX_LAYERED,
        CLASS_NAME,
        "Schublade",
        WS_POPUP,

        CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT,

        NULL,
        NULL,
        hInstance,
        NULL
    );

    if (hwnd == NULL)
    {
        return 0;
    }

    SetLayeredWindowAttributes(hwnd, 0, (255 * OPACITY) / 100, LWA_ALPHA);

    ShowWindow(hwnd, nCmdShow);
    UpdateWindow(hwnd);

    std::cout << "Created Window!" << std::endl;

    MSG msg = {};
    while (GetMessage(&msg, NULL, 0, 0) > 0)
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    UnregisterClass(CLASS_NAME, hInstance);

    return 0;
}
