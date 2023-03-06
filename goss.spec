Name: hpe-csm-goss-package
License: MIT License
Summary: Goss is a YAML based serverspec alternative tool for validating a servers configuration.
Version: %(echo ${SPEC_VERSION})
Release: %(echo ${SPEC_RELEASE})
Vendor: Hewlett Packard Enterprise Development LP
Provides: goss
%description
Installs the Goss binary onto a Linux system.

%install
pwd
mkdir -pv ${RPM_BUILD_ROOT}/usr/bin/
cp -pv ../goss-linux-amd64 ${RPM_BUILD_ROOT}/usr/bin/goss

%files
%license ../../LICENSE
%defattr(755,root,root)
/usr/bin/goss
