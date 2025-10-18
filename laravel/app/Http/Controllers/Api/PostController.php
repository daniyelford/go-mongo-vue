<?php
namespace App\Http\Controllers\Api;
use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;
class PostController extends GoApiController
{
    protected $routes = [
        'getAll' => ['url'=>'/api/posts/all','method'=>'POST'],
        'create' => ['url'=>'/api/posts/create','method'=>'POST'],
        'edit' => ['url'=>'/api/posts/edit','method'=>'PUT'],
        'delete' => ['url'=>'/api/posts/delete','method'=>'DELETE'],
    ];
    public function call(Request $request, $action)
    {
        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }
        $endpoint = $this->routes[$action]['url'];
        $method = $this->routes[$action]['method'];
        return $this->sendToGo($endpoint, $request->all(),$method,$request);
    }
}
