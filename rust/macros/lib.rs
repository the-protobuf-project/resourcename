//! Procedural derive macro powering `resourcename::Resource`.
//!
//! The [`Resource`] derive reads a single `#[resource_name(template = "...")]`
//! attribute and generates type-level `resource_template`, `parse`, and
//! `generate` helpers backed by [`resourcename::ResourceTemplate`].

use proc_macro::TokenStream;
use quote::quote;
use syn::{Attribute, DeriveInput, LitStr, parse_macro_input};

/// Derives parse/generate helpers from a `#[resource_name(template = "...")]` attribute.
///
/// Exactly one `#[resource_name(...)]` attribute is required; declaring it more
/// than once is a compile error.
#[proc_macro_derive(Resource, attributes(resource_name))]
pub fn derive_resource(input: TokenStream) -> TokenStream {
    let input = parse_macro_input!(input as DeriveInput);
    match expand(&input) {
        Ok(tokens) => tokens.into(),
        Err(err) => err.to_compile_error().into(),
    }
}

fn expand(input: &DeriveInput) -> syn::Result<proc_macro2::TokenStream> {
    let name = &input.ident;

    let attrs: Vec<&Attribute> = input
        .attrs
        .iter()
        .filter(|attr| attr.path().is_ident("resource_name"))
        .collect();

    let template = match attrs.as_slice() {
        [] => {
            return Err(syn::Error::new_spanned(
                name,
                "missing `#[resource_name(template = \"...\")]` attribute",
            ));
        }
        [attr] => parse_template(attr)?,
        [_, second, ..] => {
            return Err(syn::Error::new_spanned(
                second,
                "duplicate `resource_name` attribute: expected exactly one",
            ));
        }
    };

    let (impl_generics, ty_generics, where_clause) = input.generics.split_for_impl();

    Ok(quote! {
        impl #impl_generics #name #ty_generics #where_clause {
            /// Returns the compiled resource template for this type.
            pub fn resource_template()
                -> ::core::result::Result<
                    ::resourcename::ResourceTemplate,
                    ::resourcename::ResourceNameError,
                >
            {
                ::resourcename::ResourceTemplate::new(#template)
            }

            /// Parses a resource name into this type using its template.
            pub fn parse(resource_name: &str)
                -> ::core::result::Result<Self, ::resourcename::ResourceNameError>
            where
                Self: ::serde::de::DeserializeOwned,
            {
                Self::resource_template()?.parse_into(resource_name)
            }

            /// Generates a resource name from this value using its template.
            pub fn generate(&self)
                -> ::core::result::Result<::std::string::String, ::resourcename::ResourceNameError>
            where
                Self: ::serde::Serialize,
            {
                Self::resource_template()?.generate_from(self)
            }
        }
    })
}

fn parse_template(attr: &Attribute) -> syn::Result<LitStr> {
    let mut template: Option<LitStr> = None;
    attr.parse_nested_meta(|meta| {
        if meta.path.is_ident("template") {
            template = Some(meta.value()?.parse()?);
            Ok(())
        } else {
            Err(meta.error("expected `template = \"...\"`"))
        }
    })?;
    template.ok_or_else(|| {
        syn::Error::new_spanned(attr, "expected `#[resource_name(template = \"...\")]`")
    })
}
