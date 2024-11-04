#include <iostream>
#include <windows.h>
#include <d2d1.h>
#include <dwrite.h>

#define OPACITY 100
#define WIDTH_PERCENTAGE 0.8
#define HEIGHT_PERCENTAGE 0.5
#define DROPDOWN_DISTANCE_IN_PIXEL_PER_TICK 20
#define DROPDOWN_DURATION_IN_MILLISECONDS 200

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
    static ID2D1Factory* d2d_factory = nullptr;
    static IDWriteFactory* dwrite_factory = nullptr;
    static ID2D1HwndRenderTarget* render_target = nullptr;

    static bool isAnimating = false;

    switch (uMsg)
    {
        case WM_CREATE: {
            CoInitialize(NULL);
            D2D1CreateFactory(D2D1_FACTORY_TYPE_SINGLE_THREADED, &d2d_factory);
            DWriteCreateFactory(DWRITE_FACTORY_TYPE_SHARED, __uuidof(IDWriteFactory), reinterpret_cast<IUnknown**>(&dwrite_factory));
            return 0;
        }
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
            BeginPaint(hwnd, &ps);

            if (!render_target) {
                RECT rc;
                GetClientRect(hwnd, &rc);
                D2D1_SIZE_U size = D2D1::SizeU(rc.right - rc.left, rc.bottom - rc.top);
                d2d_factory->CreateHwndRenderTarget(
                    D2D1::RenderTargetProperties(),
                    D2D1::HwndRenderTargetProperties(hwnd, size),
                    &render_target
                );
            }

            render_target->BeginDraw();
            render_target->Clear(D2D1::ColorF(0.11f, 0.09f, 0.15f));

            ID2D1SolidColorBrush* brush = nullptr;
            render_target->CreateSolidColorBrush(D2D1::ColorF(D2D1::ColorF::White), &brush);

            IDWriteTextFormat* text_format = nullptr;
            dwrite_factory->CreateTextFormat(
                L"Agave Nerd Font",
                nullptr,
                DWRITE_FONT_WEIGHT_NORMAL,
                DWRITE_FONT_STYLE_NORMAL,
                DWRITE_FONT_STRETCH_NORMAL,
                14.0f,
                L"en-us",
                &text_format
            );

            text_format->SetTextAlignment(DWRITE_TEXT_ALIGNMENT_CENTER);
            text_format->SetParagraphAlignment(DWRITE_PARAGRAPH_ALIGNMENT_CENTER);

            RECT rc;
            GetClientRect(hwnd, &rc);
            D2D1_RECT_F layout_rect = D2D1::RectF(static_cast<float>(rc.left), static_cast<float>(rc.top), static_cast<float>(rc.right), static_cast<float>(rc.bottom));

            render_target->DrawText(
                L"Welcome to the Terminal Emulator!",
                wcslen(L"Welcome to the Terminal Emulator!"),
                text_format,
                layout_rect,
                brush
            );

            render_target->EndDraw();

            brush->Release();
            text_format->Release();
            EndPaint(hwnd, &ps);
            return 0;
        }
        case WM_DESTROY: {
            if (render_target) {
                render_target->Release();
                render_target = nullptr;
            }
            if (d2d_factory) {
                d2d_factory->Release();
                d2d_factory = nullptr;
            }
            if (dwrite_factory) {
                dwrite_factory->Release();
                dwrite_factory = nullptr;
            }
            CoUninitialize();
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
    CoInitialize(NULL);

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
