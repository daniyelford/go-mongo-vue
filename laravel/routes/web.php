<?php
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\UserController;
use App\Http\Controllers\Api\PostController;
Route::get('/health', function () {
    return response()->json(['status' => 'OK'], 200);
});
Route::post('/login/fingerPrint/has', [AuthController::class, 'call'])->defaults('action', 'hasFingerPrint');
Route::post('/login/send', [AuthController::class, 'call'])->defaults('action', 'sendCode');
Route::post('/login/verify', [AuthController::class, 'call'])->defaults('action', 'verifyCode');
Route::post('/login/fingerPrint/start', [AuthController::class, 'call'])->defaults('action', 'loginFingerPrintStart');
Route::post('/login/fingerPrint/end', [AuthController::class, 'call'])->defaults('action', 'loginFingerPrintEnd');
Route::post('/token/refresh', [AuthController::class, 'call'])->defaults('action', 'refreshToken');

// Private routes
Route::middleware(['jwt.auth'])->group(function () {
    Route::post('/register/fingerPrint/start', [AuthController::class, 'call'])->defaults('action', 'registerFingerPrintStart');
    Route::post('/register/fingerPrint/end', [AuthController::class, 'call'])->defaults('action', 'registerFingerPrintEnd');
    Route::get('/auth/logout', [AuthController::class, 'call'])->defaults('action', 'logout');
    Route::get('/auth/validate', [AuthController::class, 'call'])->defaults('action', 'validateToken');
    Route::post('/register/save', [AuthController::class, 'call'])->defaults('action', 'register');

    Route::put('/user/update', [UserController::class, 'call'])->defaults('action', 'update');
    Route::get('/user/info', [UserController::class, 'call'])->defaults('action', 'info');

    Route::post('/posts/all', [PostController::class, 'call'])->defaults('action', 'getAll');
    Route::post('/posts/create', [PostController::class, 'call'])->defaults('action', 'create');
    Route::put('/posts/edit', [PostController::class, 'call'])->defaults('action', 'edit');
    Route::delete('/posts/delete', [PostController::class, 'call'])->defaults('action', 'delete');
});