{ lib, buildGoModule }:
buildGoModule rec {
  pname = "hochwassernachrichtendienst_exporter";
  version = "1.0.0";

  src = ./.;

  vendorSha256 = "sha256-MlBoieRcJ6ggiz7NJ5DedIsWOVwG7T1GptNh43aQ01I";
  proxyVendor = true;

  meta = with lib; {
    description = "an unofficial prometheus exporter for the Hochwassernachrichtendienst Bayern";
    homepage = "https://github.com/mguentner/hochwassernachrichtendienst_exporter";
    license = licenses.mit;
    maintainers = with maintainers; [ mguentner ];
  };
}
