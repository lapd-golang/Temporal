package io

import (
	"context"
	"fmt"
	"testing"

	mdtest "gx/ipfs/QmPJNbVw8o3ohC43ppSXyNXwYKsWShG4zygnirHptfbHri/go-merkledag/test"
	ipld "gx/ipfs/QmZ6nzCLwGLVfRzYLpD7pW6UNuBDKEcA2imJtVpbEx2rxy/go-ipld-format"

	ft "gx/ipfs/QmcYUTQ7tBZeH1CLsZM2S3xhMEZdvUgXvbjhpMsLDpk3oJ/go-unixfs"
)

func TestEmptyNode(t *testing.T) {
	n := ft.EmptyDirNode()
	if len(n.Links()) != 0 {
		t.Fatal("empty node should have 0 links")
	}
}

func TestDirectoryGrowth(t *testing.T) {
	ds := mdtest.Mock()
	dir := NewDirectory(ds)
	ctx := context.Background()

	d := ft.EmptyDirNode()
	ds.Add(ctx, d)

	nelems := 10000

	for i := 0; i < nelems; i++ {
		err := dir.AddChild(ctx, fmt.Sprintf("dir%d", i), d)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := dir.GetNode()
	if err != nil {
		t.Fatal(err)
	}

	links, err := dir.Links(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(links) != nelems {
		t.Fatal("didnt get right number of elements")
	}

	dirc := d.Cid()

	names := make(map[string]bool)
	for _, l := range links {
		names[l.Name] = true
		if !l.Cid.Equals(dirc) {
			t.Fatal("link wasnt correct")
		}
	}

	for i := 0; i < nelems; i++ {
		dn := fmt.Sprintf("dir%d", i)
		if !names[dn] {
			t.Fatal("didnt find directory: ", dn)
		}

		_, err := dir.Find(context.Background(), dn)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDuplicateAddDir(t *testing.T) {
	ds := mdtest.Mock()
	dir := NewDirectory(ds)
	ctx := context.Background()
	nd := ft.EmptyDirNode()

	err := dir.AddChild(ctx, "test", nd)
	if err != nil {
		t.Fatal(err)
	}

	err = dir.AddChild(ctx, "test", nd)
	if err != nil {
		t.Fatal(err)
	}

	lnks, err := dir.Links(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(lnks) != 1 {
		t.Fatal("expected only one link")
	}
}

func TestDirBuilder(t *testing.T) {
	ds := mdtest.Mock()
	dir := NewDirectory(ds)
	ctx := context.Background()

	child := ft.EmptyDirNode()
	err := ds.Add(ctx, child)
	if err != nil {
		t.Fatal(err)
	}

	count := 5000

	for i := 0; i < count; i++ {
		err := dir.AddChild(ctx, fmt.Sprintf("entry %d", i), child)
		if err != nil {
			t.Fatal(err)
		}
	}

	dirnd, err := dir.GetNode()
	if err != nil {
		t.Fatal(err)
	}

	links, err := dir.Links(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(links) != count {
		t.Fatal("not enough links dawg", len(links), count)
	}

	adir, err := NewDirectoryFromNode(ds, dirnd)
	if err != nil {
		t.Fatal(err)
	}

	links, err = adir.Links(ctx)
	if err != nil {
		t.Fatal(err)
	}

	names := make(map[string]bool)
	for _, lnk := range links {
		names[lnk.Name] = true
	}

	for i := 0; i < count; i++ {
		n := fmt.Sprintf("entry %d", i)
		if !names[n] {
			t.Fatal("COULDNT FIND: ", n)
		}
	}

	if len(links) != count {
		t.Fatal("wrong number of links", len(links), count)
	}

	linkResults := dir.EnumLinksAsync(ctx)

	asyncNames := make(map[string]bool)
	var asyncLinks []*ipld.Link

	for linkResult := range linkResults {
		if linkResult.Err != nil {
			t.Fatal(linkResult.Err)
		}
		asyncNames[linkResult.Link.Name] = true
		asyncLinks = append(asyncLinks, linkResult.Link)
	}

	for i := 0; i < count; i++ {
		n := fmt.Sprintf("entry %d", i)
		if !asyncNames[n] {
			t.Fatal("COULDNT FIND: ", n)
		}
	}

	if len(asyncLinks) != count {
		t.Fatal("wrong number of links", len(asyncLinks), count)
	}
}
