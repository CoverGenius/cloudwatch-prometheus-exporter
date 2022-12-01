# SPEC file overview:
# https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/#con_rpm-spec-file-overview
# Fedora packaging guidelines:
# https://docs.fedoraproject.org/en-US/packaging-guidelines/


Name: cloudwatch-prometheus-exporter
Version: 0.1.8
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
* Thu Dec 1 2022 Serghei Anicheev <serghei@covergenius.com>
- Sanitize tag values used in custom checks
- Sorting tag names before building resource descriptor
* Sat Oct 29 2022 Serghei Anicheev <serghei@covergenius.com>
- Add AWS/Backup namespace
* Mon Jan 10 2022 Andrew Wright <andrew.w@covergenius.com>
- Add database memory usage metric for elasticache
* Tue Dec 14 2021 Serghei Anicheev <serghei@covergenius.com>
- Expose IP address usage per subnet
* Wed Oct 21 2020 Serghei Anicheev <serghei@covergenius.com>
- Expose existing AWS tags as a separate label
* Tue Oct 20 2020 Andrew Wright <andrew.w@covergenius.com>
- Remove default period from most metrics
* Tue Jun 02 2020 Serghei Anicheev <serghei@covergenius.com>
- Adding SQS metrics
* Thu Apr 16 2020 Andrew Wright <andrew.w@covergenius.com>
- Fix race condition when fetching resource list
* Wed Apr 15 2020 Andrew Wright <andrew.w@covergenius.com>
- Allow runtime configuration of metrics
* Wed Apr 08 2020 Andrew Wright <andrew.w@covergenius.com>
- Don't persist gauge metrics in between scrapes
* Mon Apr 06 2020 Andrew Wright <andrew.w@covergenius.com>
- Don't wait poll interval for the first scrape
* Tue Mar 31 2020 Andrew Wright <andrew.w@covergenius.com>
- Automatically determine the metric type based on the metric aggregation used
* Mon Mar 30 2020 Andrew Wright <andrew.w@covergenius.com>
- Use the official prometheus client library to manage internal metric representation
* Wed Mar 25 2020 Andrew Wright <andrew.w@covergenius.com>
- Add the ability to pick one or more different statistics for each metric
* Mon Mar 23 2020 Andrew Wright <andrew.w@covergenius.com>
- Don't publish any data for a metric if cloudwatch returns an empty results set
* Fri Mar 20 2020 Andrew Wright <andrew.w@covergenius.com>
- Use the value of the "Name" tag to populate the ec2 name label if preset
* Tue Feb 11 2020 Serghei Anicheev <serghei@covergenius.com>
- Now can specify length in config.yaml
* Tue Dec 17 2019 Serghei Anicheev <serghei@covergenius.com>
- Several fixes: map protection from concurrent access and DescribeTags limit
* Tue Nov 19 2019 Serghei Anicheev <serghei@covergenius.com>
- Initial commit
