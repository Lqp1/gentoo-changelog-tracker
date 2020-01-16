# Copyright 2020 Gentoo Authors
# Distributed under the terms of the GNU General Public License v2

EAPI=7

DESCRIPTION="Simple Live EBUILD for Gentoo-changelog-tracker"
HOMEPAGE="https://github.com/lqp1/gentoo-changelog-tracker"
EGO_PN="github.com/lqp1/gentoo-changelog-tracker"

inherit golang-vcs golang-build

KEYWORDS="~amd64"
LICENSE="GPL-3"
SLOT="0"
IUSE=""

DEPEND="app-portage/gentoolkit"
RDEPEND="${DEPEND}"

src_prepare() {
	# Ugly hack to bring all dependencies when building go sources
	cp -vr "${EGO_STORE_DIR}/src/" "${WORKDIR}/${P}/"
	default
}

src_install(){
	# Not sure why this is not detected properly ?
	dobin gentoo-changelog-tracker
}
