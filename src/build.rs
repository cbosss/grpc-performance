// fn main() {
//     let mut config = prost_build::Config::new();
//     config.type_attribute(".", "#[derive(Hash, Eq)]");
//
//     tonic_build::configure()
//         .build_client(true)
//         .build_server(true)
//         .compile_with_config(config, &["proto/echo.proto"], &[])
//         .expect("failed to compile protos");
// }
