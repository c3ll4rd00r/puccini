#!/usr/bin/env python

# Supported environment variables:
#
# PUCCINI_REPO: set this to override the Puccini git repo location (i.e. to use a local clone) 
# PUCCINI_GO_VERSION: set this to override the Go distribution version used to compile Puccini

import os, os.path, setuptools, subprocess, tempfile, shutil

with open(os.path.join(os.path.dirname(__file__), 'README.rst')) as f:
    readme = f.read()

version = os.environ.get('PUCCINI_VERSION', '0.0.0')
repo = os.environ.get('PUCCINI_REPO', 'https://github.com/tliron/puccini')
go_version = os.environ.get('PUCCINI_GO_VERSION', '1.15.7')
root = os.path.abspath(os.path.dirname(__file__)).replace('"', '\\"')

script = '''\
ROOT="{root}"

# Install Go
curl https://storage.googleapis.com/golang/go{go_version}.linux-amd64.tar.gz --silent --location | tar -xz
export PATH="$PATH:go/bin"

# Fetch repository
REPO=puccini
git clone {repo} "$REPO"

# Build library
"$REPO/scripts/build-library"
mv "$REPO/dist/libpuccini.so" "$ROOT/puccini/" 
'''.format(root=root, go_version=go_version, repo=repo)

t = tempfile.mkdtemp()
try:
    subprocess.check_output(('bash',  '-o', 'pipefail', '-euxc', script), cwd=t)
finally:
    shutil.rmtree(t)

setuptools.setup(
    name='puccini',
    version=version,
    description='Puccini',
    long_description=readme,
    license='Apache License 2.0',
    author='Tal Liron',
    author_email='tal.liron@gmail.com',
    url='https://github.com/tliron/puccini',
    download_url='https://github.com/tliron/puccini/releases',
    classifiers=[
        'Development Status :: 4 - Beta',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License'],

    packages=['puccini'],
    package_dir={'puccini': 'puccini'},
    package_data={'puccini': ['libpuccini.so']},
    install_requires=['ruamel.yaml'])
