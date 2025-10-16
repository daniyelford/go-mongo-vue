<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;

class AuthController extends GoApiController
{
    // آرایه endpoint ها: key = نام متد، value = endpoint Go
    protected $routes = [
        'hasFingerPrint' => '/auth/has-finger-print',
        'sendCode' => '/auth/send-code',
        'verifyCode' => '/auth/verify-code',
        'loginFingerPrintStart' => '/auth/login-finger-start',
        'loginFingerPrintEnd' => '/auth/login-finger-end',
        'refreshToken' => '/auth/refresh-token',
        'logout' => '/auth/logout',
        'validateToken' => '/auth/validate-token',
        'register' => '/auth/register',
        'registerFingerPrintStart' => '/auth/register-finger-start',
        'registerFingerPrintEnd' => '/auth/register-finger-end',
    ];

    // فانکشن عمومی که همه درخواست‌ها ازش عبور می‌کنن
    public function call(Request $request, $action)
    {
        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }

        $endpoint = $this->routes[$action];
        return $this->sendToGo($endpoint, $request->all());
    }
}
