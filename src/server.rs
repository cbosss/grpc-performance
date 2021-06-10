use echo::echoer_server::{Echoer, EchoerServer};
use echo::{EchoRequest, EchoResponse};
use tonic::{transport::Server, Request, Response, Status};

pub mod echo {
    tonic::include_proto!("echo"); // The string specified here must match the proto package name

    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("echo_descriptor");
}

#[derive(Debug, Default)]
pub struct EchoServer {}

#[tonic::async_trait]
impl Echoer for EchoServer {
    async fn echo(&self, request: Request<EchoRequest>) -> Result<Response<EchoResponse>, Status> {
        let resp = echo::EchoResponse {
            msg: request.into_inner().msg,
        };
        Ok(Response::new(resp))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let reflection = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(echo::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();

    let server = EchoServer::default();

    Server::builder()
        .add_service(EchoerServer::new(server))
        .add_service(reflection)
        .serve("0.0.0.0:5555".parse()?)
        .await?;

    Ok(())
}
