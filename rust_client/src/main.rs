use studentgrpc::student_client::StudentClient;
use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use studentgrpc::StudentRequest;
use serde::{Deserialize, Serialize};

pub mod studentgrpc {
    tonic::include_proto!("studentgrpc");
}

#[derive(Deserialize, Serialize)]
struct StudentData {
    name: String,
    age: i32,
    faculty: String,
    discipline: i32,
}


async fn handle_student(student: web::Json<StudentData>) -> impl Responder {
    // Definir las URLs según la disciplina
    let server_address = match student.discipline {
        1 => "http://go-server-service-swimming:50051", // swimming
        2 => "http://go-server-service-athletics:50052", // athletics
        3 => "http://go-server-service-boxing:50053", // boxing
        _ => return HttpResponse::BadRequest().body("Invalid discipline"),
    };

    let mut client = match StudentClient::connect(server_address).await {
        Ok(client) => client,
        Err(e) => return HttpResponse::InternalServerError().body(format!("Failed to connect to gRPC server: {}", e)),
    };

    let request = tonic::Request::new(StudentRequest {
        name: student.name.clone(),
        age: student.age,
        faculty: student.faculty.clone(),
        discipline: student.discipline,
    });

    match client.get_student(request).await {
            Ok(response) => {
                println!("RESPONSE={:?}", response);
                // Retornar un JSON con message vacío
                HttpResponse::Ok().body(r#"{"message": ""}"#) // Devuelve un mensaje vacío
            },
            Err(e) => HttpResponse::InternalServerError().body(format!("gRPC call failed: {}", e)),
        }
}


#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server rust at http://localhost:8080");
    HttpServer::new(|| {
        App::new()
            .route("/ingenieria", web::post().to(handle_student))
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
