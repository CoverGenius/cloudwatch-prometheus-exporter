Name: cloudwatch-prometheus-exporter
Version: 0.0.2
Release: 0%{?dist}
Summary: Cloudwatch Prometheus Exporter
License: BSD
URL: https://www.covergenius.com
BuildRequires: golang >= 1.12
BuildArch: x86_64


%description
Cloudwatch Prometheus Exporter


%prep
mkdir -p %{_topdir}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
cp -rf %{_sourcedir}/* %{_topdir}/BUILD


%build
go build


%install
rm -rf %{buildroot}
%{__install} -D -m 0644 %{_topdir}/BUILD/cloudwatch-prometheus-exporter %{buildroot}/%{_sbindir}/cloudwatch-prometheus-exporter


%files
%defattr(755,root,root,755)
%{_sbindir}/cloudwatch-prometheus-exporter


%changelog
* Tue Nov 19 2019 Serghei Anicheev <serghei@covergenius.com>
- Initial commit
